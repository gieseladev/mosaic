import abc
from typing import Any, Iterable, Mapping, Tuple

from PIL import Image

__all__ = ["CompositionABC", "Composition"]


class CompositionABC(abc.ABC):

    @abc.abstractmethod
    def draw_to(self, canvas: Image.Image) -> None:
        """Draw the composition to the given image.

        Args:
            canvas: Image to draw to.
        """
        ...


class Composition(CompositionABC, abc.ABC):
    """Like `CompositionABC` but with some utilities.


    """

    _images: Tuple[Image.Image]
    _args: Mapping[str, Any]

    def __init__(self, images: Iterable[Image.Image], **kwargs: Any) -> None:
        self._images = tuple(images)
        self._args = kwargs
