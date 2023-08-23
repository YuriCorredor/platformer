import pygame
from typing import Tuple, List

NEIGHBOURS_OFFSET = [
  (-1, 0), (-1, 1), (0, -1),
  (1, -1), (1, 0), (1, 1),
  (0, 0), (-1, 1), (0, 1),
  (1, 0), (1, 1), (0, 1)
]
PHYSICS_TILES = {'grass', 'stone'}

class Tilemap:
  def __init__(self, game, tile_size: int=16) -> None:
    self.game = game
    self.tile_size = tile_size
    self.tilemap = {}
    self.offgrid_tiles = []

    for i in range(10):
      self.tilemap[str(3 + i) + ';10'] = { 'type': 'grass', 'variant': 1, 'pos': (3 + i, 10) }
      self.tilemap['10' + ';' + str(5 + i)] = { 'type': 'stone', 'variant': 1, 'pos': (10, 5 + i) }

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
    