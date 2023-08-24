import os
import pygame
from typing import List

BASE_IMG_PATH = 'data/images/'

def load_image(path: str) -> pygame.Surface:
  image = pygame.image.load(BASE_IMG_PATH + path).convert()
  image.set_colorkey((0, 0, 0)) # black will become transparent
  return image

def load_images(path: str) -> List[pygame.Surface]:
  images = []
  for img_name in sorted(os.listdir(BASE_IMG_PATH + path)):
    images.append(load_image(path + '/' + img_name))
  return images

class Animation:
  def __init__(self, images: List[pygame.Surface], image_duration: int=5, loop: bool=True) -> None:
    self.images = images
    self.image_duration = image_duration
    self.loop = loop
    self.done = False
    self.frame = 0

  def copy(self) -> 'Animation':
    return Animation(self.images, self.image_duration, self.loop)
  
  def image(self) -> pygame.Surface:
    return self.images[int(self.frame / self.image_duration)]
  
  def update(self) -> None:
    if self.loop:
      self.frame = (self.frame + 1) % (len(self.images) * self.image_duration)
    else:
      self.frame = min(self.frame + 1, (len(self.images) - 1) * self.image_duration)
      if self.frame >= self.image_duration * (len(self.images) - 1):
        self.done = True