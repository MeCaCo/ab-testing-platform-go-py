from pydantic import BaseModel
from datetime import datetime
from typing import Optional

class TestCreate(BaseModel):
    name: str
    description: Optional[str] = None
    traffic_split: int = 50

class TestUpdate(BaseModel):
    name: Optional[str] = None
    description: Optional[str] = None
    traffic_split: Optional[int] = None
    status: Optional[str] = None

class TestResponse(BaseModel):
    id: str
    name: str
    description: Optional[str]
    status: str
    traffic_split: int
    created_at: datetime
    
    class Config:
        from_attributes = True