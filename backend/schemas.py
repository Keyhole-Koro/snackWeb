from pydantic import BaseModel
from typing import List, Optional, Dict, Any

class FeedItem(BaseModel):
    timestamp: str
    event_type: str
    agent_name: str
    content: str
    related_to: Optional[str] = None
    metadata: Dict[str, Any] = {}

class PersonaStats(BaseModel):
    post_quality: float
    incisiveness: float
    judiciousness: float
    raw_fitness: float

class PersonaDetail(BaseModel):
    name: str
    bio: str
    stats: Optional[PersonaStats] = None

class GenerationStats(BaseModel):
    generation: int
    population_diversity: float
    fitness_mean: float
