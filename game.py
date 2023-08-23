import sys
import pygame

from scripts.utils import load_image, load_images, Animation
from scripts.entities import PhysicsEntity, Player
from scripts.tilemap import Tilemap
from scripts.clouds import Clouds

PLAYER_IMG_PATH = 'entities/player.png'
DECOR_IMGS_PATH = 'tiles/decor'
GRASS_IMGS_PATH = 'tiles/grass'
LARGE_DECOR_IMGS_PATH = 'tiles/large_decor'
SPAWNERS_IMGS_PATH = 'tiles/spawners'
STONE_IMGS_PATH = 'tiles/stone'
BACKGROUND_IMG_PATH = 'background.png'
CLOUDS_IMGS_PATH = 'clouds'
PLAYER_IDLE_IMGS_PATH = 'entities/player/idle'
PLAYER_JUMP_IMGS_PATH = 'entities/player/jump'
PLAYER_RUN_IMGS_PATH = 'entities/player/run'
PLAYER_SLIDE_IMGS_PATH = 'entities/player/slide'
PLAYER_WALL_SLIDE_IMGS_PATH = 'entities/player/wall_slide'

class Game:
  def __init__(self) -> None:
    pygame.init()
    pygame.display.set_caption("Platformer")

    self.screen = pygame.display.set_mode((640, 480))
    self.display = pygame.Surface((320, 240)) # half the resolution of the screen
    self.clock = pygame.time.Clock()

    self.movement = [False, False]
    self.assets = {
      'player': load_image(PLAYER_IMG_PATH),
      'decor': load_images(DECOR_IMGS_PATH),
      'grass': load_images(GRASS_IMGS_PATH),
      'large_decor': load_images(LARGE_DECOR_IMGS_PATH),
      'spawners': load_images(SPAWNERS_IMGS_PATH),
      'stone': load_images(STONE_IMGS_PATH),
      'background': load_image(BACKGROUND_IMG_PATH),
      'clouds': load_images(CLOUDS_IMGS_PATH),
      'player_idle': Animation(load_images(PLAYER_IDLE_IMGS_PATH), image_duration=6),
      'player_run': Animation(load_images(PLAYER_RUN_IMGS_PATH), image_duration=4),
      'player_jump': Animation(load_images(PLAYER_JUMP_IMGS_PATH)),
      'player_slide': Animation(load_images(PLAYER_SLIDE_IMGS_PATH)),
      'player_wall_slide': Animation(load_images(PLAYER_WALL_SLIDE_IMGS_PATH)),
    }

    self.clouds = Clouds(self.assets['clouds'], count=16)

    self.player = Player(self, (50, 50), (8, 15))

    self.tilemap = Tilemap(self)

    self.scroll = [0, 0]

  def run(self):
    while True:
      self.display.blit(self.assets['background'], (0, 0))

      self.scroll[0] += (self.player.rect().centerx - self.display.get_width() / 2 - self.scroll[0]) / 30
      self.scroll[1] += (self.player.rect().centery - self.display.get_height() / 2 - self.scroll[1]) / 30
      render_scroll = (int(self.scroll[0]), int(self.scroll[1]))

      self.clouds.update()
      self.clouds.render(self.display, offset=render_scroll)

      self.tilemap.render(self.display, offset=render_scroll)

      self.player.update(self.tilemap, (self.movement[1] - self.movement[0], 0))
      self.player.render(self.display, offset=render_scroll)

      for event in pygame.event.get():
        self.handle_events(event)

      self.screen.blit(
        pygame.transform.scale(self.display, self.screen.get_size()),
        (0, 0)
      )
      pygame.display.update()
      self.clock.tick(60)

  def handle_events(self, event: pygame.event.Event) -> None:
    if event.type == pygame.QUIT:
      self.handle_quit()
    if event.type == pygame.KEYDOWN:
      self.handle_key_down(event.key)
    if event.type == pygame.KEYUP:
      self.handle_key_up(event.key)

  def handle_key_down(self, key: int) -> None:
    if key == pygame.K_LEFT:
      self.movement[0] = True
    if key == pygame.K_UP:
      self.player.velocity[1] = -3
    if key == pygame.K_RIGHT:
      self.movement[1] = True

  def handle_key_up(self, key: int) -> None:
    if key == pygame.K_LEFT:
      self.movement[0] = False
    if key == pygame.K_RIGHT:
      self.movement[1] = False

  def handle_quit(self):
    pygame.quit()
    sys.exit()

game = Game()
game.run()
