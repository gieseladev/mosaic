from PIL import Image, ImageDraw

import mosaic

__all__ = ["PieChart"]


class PieChart(mosaic.Composition):
    """
    Composition which displays the images like a pie chart.
    """

    def draw_to(self, canvas: Image.Image) -> None:
        # _create_pie_chart

        if len(images) > 5:
            images = random.sample(images, 5)

        final = Image.new("RGBA", (size, size))

        angle = 360 / len(images)

        for ind, image in enumerate(images):
            mask = Image.new("1", (size, size))
            draw = ImageDraw.Draw(mask)

            draw.pieslice((0, 0, size, size), -90 + ind * angle, -90 + (ind + 1) * angle, fill=1)

            image = image.resize((size, size))
            final.paste(image, mask=mask)

        return final
