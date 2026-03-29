from fastapi import FastAPI
from app.core.database import engine, Base
from app.api import tests, events
from app.models import User, Test, Event, reports

# Создаём таблицы (для разработки)
Base.metadata.create_all(bind=engine)

app = FastAPI(title="A/B Testing Platform")

# Подключаем роутеры
app.include_router(tests.router)
app.include_router(events.router)
app.include_router(reports.router)

@app.get("/health")
def health_check():
    return {"status": "ok"}

@app.get("/")
def root():
    return {
        "message": "A/B Testing Platform API",
        "docs": "/docs",
        "endpoints": {
            "tests": "/api/tests",
            "events": "/api/events",
            "reports": "/api/reports (soon)",
            "assign": "http://localhost:8080/api/v1/assign (Go)"
        }
    }