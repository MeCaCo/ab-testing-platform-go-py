from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List

from app.core.database import get_db
from app.models.event import Event
from app.models.test import Test
from app.schemas.event import EventCreate, EventResponse

# Создаём роутер с префиксом /api/events
router = APIRouter(prefix="/api/events", tags=["events"])

@router.post("/", response_model=EventResponse)
def create_event(event: EventCreate, db: Session = Depends(get_db)):
    """
    Отправить событие (показ/клик/конверсия)
    
    Это основной эндпоинт, куда клиентский JS отправляет данные
    """
    test = db.query(Test).filter(Test.id == event.test_id).first()
    if not test:
        raise HTTPException(status_code=404, detail="Test not found")
    
    if event.variant not in ['A', 'B']:
        raise HTTPException(status_code=400, detail="Variant must be 'A' or 'B'")
    
    valid_types = ['impression', 'click', 'conversion']
    if event.event_type not in valid_types:
        raise HTTPException(status_code=400, detail=f"Event type must be one of: {valid_types}")
    

    db_event = Event(**event.model_dump())
    db.add(db_event)
    db.commit()
    db.refresh(db_event) 
    
    return db_event

@router.get("/", response_model=List[EventResponse])
def get_events(test_id: str = None, db: Session = Depends(get_db)):
    """
    Получить список событий (можно фильтровать по test_id)
    """
    query = db.query(Event)
    
  
    if test_id:
        query = query.filter(Event.test_id == test_id)
    
    return query.order_by(Event.created_at.desc()).limit(1000).all()