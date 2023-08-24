import random
import pygame
from typing import Tuple

class Cloud:
  def __init__(self, pos: Tuple[int, int], img: pygame.Surface, speed: int, depth: int) -> None:
    self.pos = list(pos)
    self.img = img
    self.speed = speed
    self.depth = depth

  def update(self) -> None:
    self.pos[0] += self.speed

  def render(self, surface: pygame.Surface, offset: Tuple[int, int]=(0, 0)) -> None:
    render_position = (self.pos[0] - offset[0] * self.depth, self.pos[1] - offset[1] * self.depth)
    surface.blit(self.img, (render_position[0] % (surface.get_width() + self.img.get_width()) - self.img.get_width(), (render_position[1] % (surface.get_height() + self.img.get_height()) - self.img.get_height())))

class Clouds:
  def __init__(self, cloud_images: list, count: int=16) -> None:
    self.cloud_images = cloud_images
    self.clouds = []

    for _ in range(count):
      self.clouds.append(Cloud(
        (random.random() * 99999, random.random() * 99999),
        random.choice(self.cloud_images),
        random.random() * 0.05 + 0.05,
        random.random() * 0.6 + 0.2
      ))

    self.clouds.sort(key=lambda cloud: cloud.depth)

  def update(self) -> None:
    for cloud in self.clouds:
      cloud.update()
  
  def render(self, surface: pygame.Surface, offset: Tuple[int, int]=(0, 0)) -> None:
    for cloud in self.clouds:
      cloud.render(surface, offset=offset)