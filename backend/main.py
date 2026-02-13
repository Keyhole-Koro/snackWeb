from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
import os

from .services import DataService
from .schemas import FeedItem, PersonaDetail, GenerationStats
from typing import List

app = FastAPI(title="SnackWeb API")

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"], # Allow all for dev
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Service
service = DataService()

# Endpoints
@app.get("/api/feed", response_model=List[FeedItem])
def get_feed(limit: int = 50):
    return service.get_feed(limit)

@app.get("/api/personas", response_model=List[PersonaDetail])
def get_personas():
    return service.get_personas()

@app.get("/api/stats", response_model=List[GenerationStats])
def get_stats():
    return service.get_generation_stats()

# Mount plots directory
PLOTS_DIR = os.getenv("PLOTS_DIR", "/home/unix/snack/persona_data/plots")
if os.path.exists(PLOTS_DIR):
    app.mount("/plots", StaticFiles(directory=PLOTS_DIR), name="plots")

@app.get("/")
def read_root():
    return {"message": "SnackWeb Backend Running"}
