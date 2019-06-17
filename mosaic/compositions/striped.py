from PIL import Image, ImageOps

import mosaic

__all__ = ["StripedVertical", "StripedVerticalStoreyed"]


class StripedVertical(mosaic.Composition):
    """
    Composition with vertical stripes of images.
    """

    def draw_to(self, canvas: Image.Image) -> None:
        # _create_striped_collage
        if len(images) < 2:
            raise ValueError("Need at least 2 images")

        if len(images) > 5:
            images = random.sample(images, 5)

        stripe_width = size // len(images)
        canvas = Image.new("RGB", (size, size))

        for i, image in enumerate(images):
            box = (i * stripe_width, 0, (i + 1) * stripe_width, size)
            image = image.crop(box=box)
            canvas.paste(image, box=box)

        return canvas


class StripedVerticalStoreyed(mosaic.Composition):
    """
    Composition similar to `StripedVertical` but a stripe contains multiple
    images.
    """

    def draw_to(self, canvas: Image.Image) -> None:
        # _create_partitioned_striped_collage

        if len(images) < 3:
            raise ValueError("Need at least 3 images")

        min_pick = 1
        max_pick = max(round(len(images) / 3), 3)

        if len(images) == 3:
            max_pick = 2

        images = list(images)
        stripes = []
        while images:
            im_count = random.randint(min(len(images), min_pick), min(len(images), max_pick))
            stripe = [images.pop() for _ in range(im_count)]
            stripes.append(stripe)

        stripe_width = size // len(stripes)
        canvas = Image.new("RGB", (size, size), color="white")

        min_height = size // 6

        for i, stripe in enumerate(stripes):
            image_count = len(stripe)
            max_height = size - image_count * min_height
            heights = []
            start = 0
            for _ in range(image_count - 1):
                height = random.randint(min_height, max_height)
                new_start = start + height
                heights.append((start, new_start))
                start = new_start
            heights.append((start, size))

            for j, image in enumerate(stripe):
                start, stop = heights[j]
                box = (i * stripe_width, start, (i + 1) * stripe_width, stop)
                image = ImageOps.fit(image, (stripe_width, stop - start))
                canvas.paste(image, box=box)

        return canvas
