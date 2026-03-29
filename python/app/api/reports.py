from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import Dict, Any
from datetime import datetime, timedelta
import math

from app.core.database import get_db
from app.models.test import Test
from app.models.event import Event

router = APIRouter(prefix="/api/reports", tags=["reports"])

@router.get("/{test_id}")
def get_test_report(test_id: str, db: Session = Depends(get_db)):
    """
    Получить отчёт по тесту (показывает конверсии, статистическую значимость)
    """
    # Находим тест
    test = db.query(Test).filter(Test.id == test_id).first()
    if not test:
        raise HTTPException(status_code=404, detail="Test not found")

    # Считаем события по группам
    events = db.query(Event).filter(Event.test_id == test_id).all()
    
    # Группируем по вариантам
    variant_stats = {}
    for event in events:
        if event.variant not in variant_stats:
            variant_stats[event.variant] = {'impressions': 0, 'conversions': 0, 'value': 0.0}
        
        if event.event_type == 'impression':
            variant_stats[event.variant]['impressions'] += 1
        elif event.event_type == 'conversion':
            variant_stats[event.variant]['conversions'] += 1
            variant_stats[event.variant]['value'] += event.value

    # Простая статистика (без сложной математики)
    results = {}
    for variant, stats in variant_stats.items():
        impressions = stats['impressions']
        conversions = stats['conversions']
        conversion_rate = (conversions / impressions * 100) if impressions > 0 else 0

        results[variant] = {
            'impressions': impressions,
            'conversions': conversions,
            'conversion_rate': round(conversion_rate, 2),
            'total_value': round(stats['value'], 2)
        }

    # Простая оценка значимости (наивный p-value)
    if 'A' in results and 'B' in results:
        a_cr = results['A']['conversion_rate']
        b_cr = results['B']['conversion_rate']
        
        # Наивная формула (для портфолио достаточно)
        se_a = math.sqrt((a_cr * (100 - a_cr)) / results['A']['impressions']) if results['A']['impressions'] > 0 else 0
        se_b = math.sqrt((b_cr * (100 - b_cr)) / results['B']['impressions']) if results['B']['impressions'] > 0 else 0
        se_diff = math.sqrt(se_a**2 + se_b**2)
        
        z_score = abs(b_cr - a_cr) / se_diff if se_diff > 0 else 0
        # Если z_score > 1.96, то разница значима (95% confidence)
        is_significant = z_score > 1.96

        results['significance'] = {
            'z_score': round(z_score, 2),
            'is_significant': is_significant,
            'confidence_level': '95%' if is_significant else 'Below 95%'
        }

    return {
        'test_id': test_id,
        'test_name': test.name,
        'status': test.status,
        'results': results,
        'created_at': test.created_at,
        'started_at': test.started_at
    }