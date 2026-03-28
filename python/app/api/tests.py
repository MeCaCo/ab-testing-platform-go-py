from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.models.test import Test
from app.schemas.test import TestCreate, TestResponse, TestUpdate
from typing import List

router = APIRouter(prefix="/api/tests", tags=["tests"])

@router.post("/", response_model=TestResponse)
def create_test(test: TestCreate, db: Session = Depends(get_db)):
    """Создать новый А/Б тест"""
    db_test = Test(**test.dict())
    db.add(db_test)
    db.commit()
    db.refresh(db_test)
    return db_test

@router.get("/", response_model=List[TestResponse])
def get_tests(db: Session = Depends(get_db)):
    """Получить список всех тестов"""
    return db.query(Test).all()

@router.get("/{test_id}", response_model=TestResponse)
def get_test(test_id: str, db: Session = Depends(get_db)):
    """Получить тест по ID"""
    test = db.query(Test).filter(Test.id == test_id).first()
    if not test:
        raise HTTPException(status_code=404, detail="Test not found")
    return test