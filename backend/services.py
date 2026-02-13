import json
import os
from typing import List, Dict
from datetime import datetime
from .schemas import FeedItem, PersonaDetail, PersonaStats, GenerationStats

# Path setup - assuming running from root or backend dir
# We need to find the snackPersona/persona_data directory
# Heuristic: verify absolute path or relative
SIM_DATA_DIR = os.getenv("PERSONA_DATA_DIR", "/home/unix/snack/persona_data")

class DataService:
    def __init__(self):
        self.stats_path = os.path.join(SIM_DATA_DIR, "generation_stats.jsonl")
        self.events_path = os.path.join(SIM_DATA_DIR, "simulation_events.jsonl")

    def get_feed(self, limit: int = 50) -> List[FeedItem]:
        """Get latest feed items (reverse chronological)."""
        if not os.path.exists(self.events_path):
            return []
        
        items = []
        try:
            # Read file from end efficiently? For now, readlines() is fine for small-scale sim.
            with open(self.events_path, "r") as f:
                lines = f.readlines()
                
            for line in reversed(lines):
                if len(items) >= limit:
                    break
                try:
                    data = json.loads(line)
                    items.append(FeedItem(**data))
                except:
                    continue
        except Exception as e:
            print(f"Error reading feed: {e}")
            
        return items

    def get_personas(self) -> List[PersonaDetail]:
        """Get personas from the latest generation."""
        # Find latest gen file
        if not os.path.exists(SIM_DATA_DIR):
            return []
            
        gen_files = [f for f in os.listdir(SIM_DATA_DIR) if f.startswith("gen_") and f.endswith(".json")]
        if not gen_files:
            return []
            
        # Sort by gen id
        gen_files.sort(key=lambda x: int(x.split("_")[1].split(".")[0]))
        latest_file = os.path.join(SIM_DATA_DIR, gen_files[-1])
        
        personas = []
        try:
            with open(latest_file, "r") as f:
                data = json.loads(f.read())
                
            # We also want their latest stats if available in generation_stats.jsonl
            latest_stats = self._get_latest_agent_stats()
            
            for p in data:
                name = p.get("name")
                stats = latest_stats.get(name)
                personas.append(PersonaDetail(
                    name=name,
                    bio=p.get("bio", ""),
                    stats=stats
                ))
        except Exception as e:
            print(f"Error reading personas: {e}")
            
        return personas

    def _get_latest_agent_stats(self) -> Dict[str, PersonaStats]:
        """Helper to get stats map from last line of stats file."""
        if not os.path.exists(self.stats_path):
            return {}
            
        try:
            with open(self.stats_path, "r") as f:
                lines = f.readlines()
                if not lines:
                    return {}
                last_line = json.loads(lines[-1])
                
            stats_map = {}
            for agent in last_line.get("agents", []):
                stats_map[agent["name"]] = PersonaStats(
                    post_quality=agent.get("post_quality", 0),
                    incisiveness=agent.get("incisiveness", 0),
                    judiciousness=agent.get("judiciousness", 0),
                    raw_fitness=agent.get("raw_fitness", 0)
                )
            return stats_map
        except:
            return {}

    def get_generation_stats(self) -> List[GenerationStats]:
        """Get history of generation stats."""
        if not os.path.exists(self.stats_path):
            return []
            
        stats = []
        try:
            with open(self.stats_path, "r") as f:
                for line in f:
                    try:
                        d = json.loads(line)
                        stats.append(GenerationStats(
                            generation=d["generation"],
                            population_diversity=d["population_diversity"],
                            fitness_mean=d["fitness_mean"]
                        ))
                    except:
                        continue
        except:
            pass
        return stats
