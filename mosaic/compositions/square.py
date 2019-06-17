from PIL import Image

import mosaic
import math

__all__ = [
    "PerfectSquare", "FocusedSquare",
    "DiamondSquare",
]


class PerfectSquare(mosaic.Composition):
    """
    Composition dividing the canvas into x * x image squares.
    """

    def draw_to(self, canvas: Image.Image) -> None:
        # _create_normal_collage

        if len(images) < 4:
            raise ValueError("Need at least 4 images!")

        if len(images) not in (4, 9, 16):
            if len(images) > 16:
                target_length = 16
            elif len(images) > 9:
                target_length = 9
            else:
                target_length = 4
            images = random.sample(images, target_length)

        s = round(math.sqrt(len(images)))

        part_size = size // s

        final = Image.new("RGB", (size, size))

        for i in range(s):
            for j in range(s):
                left = i * part_size
                top = j * part_size

                im = images[i * s + j]

                if crop_images:
                    right = left + part_size
                    bottom = top + part_size

                    im = im.resize((size, size)).crop((left, top, right, bottom))
                else:
                    im = im.resize((part_size, part_size))

                final.paste(im, box=(left, top))

        return final


class FocusedSquare(mosaic.Composition):
    """
    Composition which prominently features a single image and fills the
    remaining space with the other images.
    """

    def draw_to(self, canvas: Image.Image) -> None:
        # _create_focused_collage

        if len(images) < 4:
            raise ValueError("Amount of images should be at least 4")

        if len(images) % 2 != 0:
            images = random.sample(images, len(images) - 1)

        final = Image.new("RGB", (size, size))

        focus_image_size = size * (len(images) - 2) // len(images)
        other_images_size = size - focus_image_size

        focus_image = images[0].resize((focus_image_size, focus_image_size))
        final.paste(focus_image, box=(0, size - focus_image_size))

        top_right = images[1].resize((other_images_size, other_images_size))
        final.paste(top_right, box=(focus_image_size, 0))

        leftover = images[2:]

        for i in range(len(leftover) // 2):
            left = leftover[i].resize((other_images_size, other_images_size))
            right = leftover[-i - 1].resize((other_images_size, other_images_size))

            final.paste(left, box=(i * other_images_size, 0))
            final.paste(right, box=(focus_image_size, (i + 1) * other_images_size))

        return final


class DiamondSquare(mosaic.Composition):

    def draw_to(self, canvas: Image.Image) -> None:
        # _create_diamond_square_collage

        if len(images) not in {1, 5, 9, 13}:
            raise ValueError("Need 1, 5, 9, or 13 images!")

        images = list(images)

        canvas: Image.Image = Image.new("RGBA", (size, size), color=None)

        inner_size: int = int(.9925 * size)
        diamond_len: int = int(3 * math.sqrt(2) * inner_size / (13 + math.sqrt(2)))

        diamond_mask = Image.new("1", (diamond_len, diamond_len), color=1).rotate(45, expand=True)

        # place center image
        img = images.pop().resize(diamond_mask.size)
        canvas.paste(img, box=get_center_pos(img, canvas), mask=diamond_mask)

        x_center: float = canvas.height / 2
        y_center: float = canvas.height / 2

        if len(images) >= 4:
            # place adjacent images
            for i_y in range(2):
                y_pos: int = int(y_center - i_y * diamond_mask.height)

                for i_x in range(2):
                    x_pos: int = int(x_center - i_x * diamond_mask.width)

                    img = images.pop().resize(diamond_mask.size)
                    canvas.paste(img, box=(x_pos, y_pos), mask=diamond_mask)

        if len(images) >= 4:
            small_dia_width: int = int(2 / 3 * diamond_mask.width)
            small_dia_mask = diamond_mask.resize((small_dia_width, small_dia_width))

            tr_s: float = small_dia_mask.width / 2
            touching_radius: float = diamond_mask.width / 2

            # place in-between small images
            for rot in range(4):
                angle: float = rot * math.pi / 2
                touching_x_pos: float = x_center + math.cos(angle) * touching_radius
                touching_y_pos: float = y_center + math.sin(angle) * touching_radius

                fac_x: int = (0, 1, 2, 1)[rot]
                x_pos: int = int(touching_x_pos - fac_x * tr_s)

                fac_y: int = (1, 0, 1, 2)[rot]
                y_pos: int = int(touching_y_pos - fac_y * tr_s)

                img = images.pop().resize(small_dia_mask.size)
                canvas.paste(img, box=(x_pos, y_pos), mask=small_dia_mask)

            if len(images) >= 4:
                dangle_offset: float = math.sqrt(2) * 1 / 4 * small_dia_width
                touching_radius: float = 1.5 * diamond_len - dangle_offset

                # place tip small images
                for rot in range(4):
                    angle: float = math.pi / 4 + rot * math.pi / 2
                    touching_x_pos: float = x_center + math.cos(angle) * touching_radius
                    touching_y_pos: float = y_center + math.sin(angle) * touching_radius

                    x_pos: int = int(touching_x_pos - (rot == 1 or rot == 2) * small_dia_width)
                    y_pos: int = int(touching_y_pos - (rot >= 2) * small_dia_width)

                    img = images.pop().resize(small_dia_mask.size)
                    canvas.paste(img, box=(x_pos, y_pos), mask=small_dia_mask)

        return canvas
