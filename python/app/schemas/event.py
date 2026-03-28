from pydantic import BaseModel
from datetime import datetime
from typing import Optional

# Базовая схема для события (какие данные мы принимаем)
class EventCreate(BaseModel):
    test_id: str         
    user_id: str          
    variant: str         
    event_type: str       
    value: float = 1.0     

# Схема ответа (что возвращаем клиенту)
class EventResponse(BaseModel):
    id: int
    test_id: str
    user_id: str
    variant: str
    event_type: str
    value: float
    created_at: datetime
    
    class Config:
        from_attributes = True