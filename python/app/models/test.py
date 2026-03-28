from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, Float
from sqlalchemy.orm import relationship
from app.core.database import Base
import datetime
import uuid

class Test(Base):
    __tablename__ = "tests"
    
    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String, index=True, nullable=False)
    description = Column(String, nullable=True)
    status = Column(String, default="draft")  # draft, active, paused, completed
    traffic_split = Column(Integer, default=50)  # % для варианта B
    owner_id = Column(Integer, ForeignKey("users.id"))
    created_at = Column(DateTime, default=datetime.datetime.utcnow)
    started_at = Column(DateTime, nullable=True)
    ended_at = Column(DateTime, nullable=True)
    
    owner = relationship("User", back_populates="tests")
    events = relationship("Event", back_populates="test")