import pygame
from typing import Tuple, List

class PhysicsEntity:
  def __init__(self, game, entity_type: str, position: Tuple[int, int], size: Tuple[int, int]) -> None:
    self.game = game
    self.type = entity_type
    self.position = list(position)
    self.size = size
    self.velocity = [0, 0]
    self.collisions = { 'up': False, 'down': False, 'left': False, 'right': False }

    self.action = ''
    self.anim_offset = (-3, -3)
    self.flip = False
    self.set_action('idle')

  def rect(self):
    return pygame.Rect(self.position[0], self.position[1], self.size[0], self.size[1])
  
  def set_action(self, action: str) -> None:
    if self.action != action:
      self.action = action
      self.anim = self.game.assets[self.type + '_' + action].copy()

  def update(self, tilemap: List[pygame.Rect], movement: Tuple[int, int]=(0, 0)) -> None:
    self.collisions = { 'up': False, 'down': False, 'left': False, 'right': False }

    frame_movement = (movement[0] + self.velocity[0], movement[1] + self.velocity[1])

    self.position[0] += frame_movement[0]
    entity_rect = self.rect()
    for rect in tilemap.physics_rects_around(self.position):
      if entity_rect.colliderect(rect):
        if frame_movement[0] > 0:
          entity_rect.right = rect.left
          self.collisions['right'] = True
        if frame_movement[0] < 0:
          entity_rect.left = rect.right
          self.collisions['left'] = True
        self.position[0] = entity_rect.x

    self.position[1] += frame_movement[1]
    entity_rect = self.rect()
    for rect in tilemap.physics_rects_around(self.position):
      if entity_rect.colliderect(rect):
        if frame_movement[1] > 0:
          entity_rect.bottom = rect.top
          self.collisions['down'] = True
        if frame_movement[1] < 0:
          entity_rect.top = rect.bottom
          self.collisions['up'] = True
        self.position[1] = entity_rect.y

    if movement[0] < 0:
      self.flip = True
    if movement[0] > 0:
      self.flip = False

    self.velocity[1] = min(5, self.velocity[1] + 0.1)

    if self.collisions['down'] or self.collisions['up']:
      self.velocity[1] = 0

    self.anim.update()

  def render(self, surface: pygame.Surface, offset: Tuple[int, int]=(0, 0)) -> None:
    surface.blit(pygame.transform.flip(self.anim.image(), self.flip, False), (self.position[0] - offset[0] + self.anim_offset[0], self.position[1] - offset[1] + self.anim_offset[1]))

class Player(PhysicsEntity):
  def __init__(self, game, position: Tuple[int, int], size: Tuple[int, int]) -> None:
    super().__init__(game, 'player', position, size)
    self.air_time = 0

  def update(self, tilemap: List[pygame.Rect], movement: Tuple[int, int] = (0, 0)) -> None:
    super().update(tilemap, movement)

    self.air_time += 1
    if self.collisions['down']:
      self.air_time = 0

    if self.air_time > 4:
      self.set_action('jump')
    elif movement[0] != 0:
      self.set_action('run')
    else:
      self.set_action('idle')

