import sys
import pygame

from scripts.utils import load_images
from scripts.tilemap import Tilemap
from typing import Tuple

RENDER_SCALE = 2.0

DECOR_IMGS_PATH = 'tiles/decor'
GRASS_IMGS_PATH = 'tiles/grass'
LARGE_DECOR_IMGS_PATH = 'tiles/large_decor'
SPAWNERS_IMGS_PATH = 'tiles/spawners'
STONE_IMGS_PATH = 'tiles/stone'

class Editor:
  def __init__(self) -> None:
    pygame.init()
    pygame.display.set_caption("Editor")

    self.screen = pygame.display.set_mode((640, 480))
    self.display = pygame.Surface((320, 240))
    self.clock = pygame.time.Clock()

    self.assets = {
      'decor': load_images(DECOR_IMGS_PATH),
      'grass': load_images(GRASS_IMGS_PATH),
      'large_decor': load_images(LARGE_DECOR_IMGS_PATH),
      'stone': load_images(STONE_IMGS_PATH),
    }

    self.movement = [False, False, False, False]

    self.tilemap = Tilemap(self)

    try:
      self.tilemap.load('map.json')
    except FileNotFoundError:
      pass

    self.tile_list = list(self.assets)
    self.tile_group = 0
    self.tile_variant = 0

    self.clicking = False
    self.right_clicking = False
    self.shift = False
    self.ongrid = True

    self.scroll = [0, 0]

  def run(self):
    while True:
      self.display.fill((0, 0, 0))

      self.scroll[0] += (self.movement[1] - self.movement[0]) * 3
      self.scroll[1] += (self.movement[3] - self.movement[2]) * 3
      render_scroll = int(self.scroll[0]), int(self.scroll[1])
      self.tilemap.render(self.display, offset=render_scroll)

      current_tile_image = self.assets[self.tile_list[self.tile_group]][self.tile_variant].copy()
      current_tile_image.set_alpha(150)
      self.display.blit(current_tile_image, (5, 5))

      mouse_pos = pygame.mouse.get_pos()
      mouse_pos = (mouse_pos[0] / RENDER_SCALE, mouse_pos[1] / RENDER_SCALE)
      tile_pos = (int(mouse_pos[0] + self.scroll[0]) // self.tilemap.tile_size, int(mouse_pos[1] + self.scroll[1]) // self.tilemap.tile_size)

      if self.ongrid:
        self.display.blit(current_tile_image, (tile_pos[0] * self.tilemap.tile_size - self.scroll[0], tile_pos[1] * self.tilemap.tile_size - self.scroll[1]))
      else:
        self.display.blit(current_tile_image, mouse_pos)

      if self.clicking and self.ongrid:
        self.tilemap.tilemap[str(tile_pos[0]) + ';' + str(tile_pos[1])] = {
          'type': self.tile_list[self.tile_group],
          'variant': self.tile_variant,
          'pos': tile_pos
        }

      if self.right_clicking:
        tile_location = str(tile_pos[0]) + ';' + str(tile_pos[1])
        if tile_location in self.tilemap.tilemap:
          del self.tilemap.tilemap[tile_location]
        for tile in self.tilemap.offgrid_tiles.copy():
          tile_img = self.assets[tile['type']][tile['variant']]
          tile_rect = pygame.Rect(tile['pos'][0] - self.scroll[0], tile['pos'][1] - self.scroll[1], tile_img.get_width(), tile_img.get_height())
          if tile_rect.collidepoint(mouse_pos):
            self.tilemap.offgrid_tiles.remove(tile)

      for event in pygame.event.get():
        if event.type == pygame.QUIT:
          self.handle_quit()

        if event.type == pygame.MOUSEBUTTONDOWN:
          self.handle_mouse_down(event.button, mouse_pos=mouse_pos, render_scroll=render_scroll)

        if event.type == pygame.MOUSEBUTTONUP:
          self.handle_mouse_up(event.button)

        if event.type == pygame.KEYDOWN:
          self.handle_key_down(event.key)

        if event.type == pygame.KEYUP:
          self.handle_key_up(event.key)

      self.screen.blit(
        pygame.transform.scale(self.display, self.screen.get_size()),
        (0, 0)
      )
      pygame.display.update()
      self.clock.tick(60)

  def handle_mouse_up(self, button: int) -> None:
    if button == 1:
      self.clicking = False

    if button == 3:
      self.right_clicking = False

  def handle_mouse_down(self, button: int, mouse_pos: Tuple[int, int]=None, render_scroll: Tuple[int, int]=None) -> None:
    if button == 1:
      self.clicking = True
      if not self.ongrid:
        self.tilemap.offgrid_tiles.append({
          'type': self.tile_list[self.tile_group],
          'variant': self.tile_variant,
          'pos': (mouse_pos[0] + render_scroll[0], mouse_pos[1] + render_scroll[1])
        })

    if button == 3:
      self.right_clicking = True

    if self.shift:
      if button == 4:
        self.tile_variant = (self.tile_variant + 1) % len(self.assets[self.tile_list[self.tile_group]]) 

      if button == 5:
        self.tile_variant = (self.tile_variant - 1) % len(self.assets[self.tile_list[self.tile_group]])
    else:
      if button == 4:
        self.tile_group = (self.tile_group + 1) % len(self.tile_list)
        self.tile_variant = 0

      if button == 5:
        self.tile_group = (self.tile_group - 1) % len(self.tile_list)
        self.tile_variant = 0

  def handle_key_down(self, key: int) -> None:
    if key == pygame.K_a:
      self.movement[0] = True

    if key == pygame.K_d:
      self.movement[1] = True

    if key == pygame.K_w:
      self.movement[2] = True

    if key == pygame.K_s:
      self.movement[3] = True

    if key == pygame.K_LSHIFT:
      self.shift = True

    if key == pygame.K_o:
      self.tilemap.save('map.json')

    if key == pygame.K_t:
      self.tilemap.auto_tile()

    if key == pygame.K_g:
      self.ongrid = not self.ongrid

  def handle_key_up(self, key: int) -> None:
    if key == pygame.K_a:
      self.movement[0] = False

    if key == pygame.K_d:
      self.movement[1] = False

    if key == pygame.K_w:
      self.movement[2] = False

    if key == pygame.K_s:
      self.movement[3] = False

    if key == pygame.K_LSHIFT:
      self.shift = False

  def handle_quit(self):
    pygame.quit()
    sys.exit()

Editor().run()

