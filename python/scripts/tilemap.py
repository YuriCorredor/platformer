import json
import pygame
from typing import Tuple

AUTOTILE_MAP = {
  tuple(sorted([(1, 0), (0, 1)])): 0,
  tuple(sorted([(1, 0), (0, 1), (-1, 0)])): 1,
  tuple(sorted([(-1, 0), (0, 1)])): 2,
  tuple(sorted([(-1, 0), (0, -1), (0, 1)])): 3,
  tuple(sorted([(-1, 0), (0, -1)])): 4,
  tuple(sorted([(-1, 0), (0, -1), (1, 0)])): 5,
  tuple(sorted([(1, 0), (0, -1)])): 6,
  tuple(sorted([(1, 0), (0, -1), (0, 1)])): 7,
  tuple(sorted([(1, 0), (-1, 0), (0, 1), (0, -1)])): 8,
}

NEIGHBOURS_OFFSET = [
  (-1, 0), (-1, 1), (0, -1),
  (1, -1), (1, 0), (1, 1),
  (0, 0), (-1, 1), (0, 1),
  (1, 0), (1, 1), (0, 1)
]
PHYSICS_TILES = {'grass', 'stone'}
AUTOTILES_TYPES = {'grass', 'stone'}

class Tilemap:
  def __init__(self, game, tile_size: int=16) -> None:
    self.game = game
    self.tile_size = tile_size
    self.tilemap = {}
    self.offgrid_tiles = []

  def save(self, path: str) -> None:
    f = open(path, 'w')
    json.dump({'tilemap': self.tilemap, 'tile_size': self.tile_size, 'offgrid_tiles': self.offgrid_tiles}, f)
    f.close()

  def auto_tile(self) -> None:
    for location in self.tilemap:
      tile = self.tilemap[location]
      neighbours = set()
      for shift in [(1, 0), (-1, 0), (0, -1), (0, 1)]:
        check_location = str(tile['pos'][0] + shift[0]) + ';' + str(tile['pos'][1] + shift[1])
        if check_location in self.tilemap:
          if self.tilemap[check_location]['type'] == tile['type']:
            neighbours.add(shift)

      neighbours = tuple(sorted(neighbours))
      if tile['type'] in AUTOTILES_TYPES and neighbours in AUTOTILE_MAP:
        tile['variant'] = AUTOTILE_MAP[neighbours]

  def load(self, path: str) -> None:
    f = open(path, 'r')
    data = json.load(f)
    self.tilemap = data['tilemap']
    self.tile_size = data['tile_size']
    self.offgrid_tiles = data['offgrid_tiles']
    f.close()

  def tiles_around(self, position: Tuple[int, int]) -> list:
    tiles = []
    tile_loc = (int(position[0] // self.tile_size), int(position[1] // self.tile_size))
    for offset in NEIGHBOURS_OFFSET:
      check_location = str(tile_loc[0] + offset[0]) + ';' + str(tile_loc[1] + offset[1])
      if check_location in self.tilemap:
        tiles.append(self.tilemap[check_location])

    return tiles
  
  def physics_rects_around(self, position: Tuple[int, int]) -> list:
    rects = []
    for tile in self.tiles_around(position):
      if tile['type'] in PHYSICS_TILES:
        rects.append(pygame.Rect(tile['pos'][0] * self.tile_size, tile['pos'][1] * self.tile_size, self.tile_size, self.tile_size))
    return rects

  def render(self, surface: pygame.Surface, offset: Tuple[int, int]=(0,0)) -> None:
    for tile in self.offgrid_tiles:
      surface.blit(self.game.assets[tile['type']][tile['variant']], (tile['pos'][0] - offset[0], tile['pos'][1] - offset[1]))

    for x in range(offset[0] // self.tile_size, (offset[0] + surface.get_width()) // self.tile_size + 1):
      for y in range(offset[1] // self.tile_size, (offset[1] + surface.get_height()) // self.tile_size + 1):
        location = str(x) + ';' + str(y)
        if location in self.tilemap:
          tile = self.tilemap[location]
          surface.blit(self.game.assets[tile['type']][tile['variant']], (tile['pos'][0] * self.tile_size - offset[0], tile['pos'][1] * self.tile_size - offset[1]))
    