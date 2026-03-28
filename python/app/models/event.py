from sqlalchemy import Column, BigInteger, String, DateTime, ForeignKey, Float
from sqlalchemy.orm import relationship
from app.core.database import Base
import datetime

class Event(Base):
    __tablename__ = "events"
    
    id = Column(BigInteger, primary_key=True, index=True)
    test_id = Column(String, ForeignKey("tests.id"), index=True, nullable=False)
    user_id = Column(String, index=True, nullable=False)  # ID пользователя на клиенте
    variant = Column(String, nullable=False)  # 'A' или 'B'
    event_type = Column(String, nullable=False)  # impression, click, conversion
    value = Column(Float, default=1.0)
    created_at = Column(DateTime, default=datetime.datetime.utcnow, index=True)
    
    test = relationship("Test", back_populates="events")