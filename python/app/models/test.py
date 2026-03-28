from sqlalchemy import Column, String, Integer, DateTime, ForeignKey, Float, Boolean
from sqlalchemy.orm import relationship
from app.core.database import Base
import datetime
import uuid

class Test(Base):
    __tablename__ = "tests"
    
    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String, index=True, nullable=False)
    description = Column(String, nullable=True)
    status = Column(String, default="draft")
    traffic_split = Column(Integer, default=50)
    owner_id = Column(Integer, ForeignKey("users.id"), nullable=True)
    created_at = Column(DateTime, default=datetime.datetime.utcnow)
    started_at = Column(DateTime, nullable=True)
    ended_at = Column(DateTime, nullable=True)
    
    # Связи
    events = relationship("Event", back_populates="test", cascade="all, delete-orphan")
    
    owner = relationship("User", back_populates="tests")  # ← ВОТ ЭТОГО НЕ ХВАТАЛО!