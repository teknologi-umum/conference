from html.parser import HTMLParser
from io import StringIO
from os import listdir
from os.path import isfile


# Stolen from https://stackoverflow.com/a/925630/3153224
# Copyright CC BY-SA 4.0 (c) 2009  Eloff
class MLStripper(HTMLParser):
    def __init__(self):
        super().__init__()
        self.reset()
        self.strict = False
        self.convert_charrefs = True
        self.text = StringIO()
        self.body_found = False

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if tag == "body":
            self.body_found = True

    def handle_endtag(self, tag: str) -> None:
        if tag == "body":
            self.body_found = False

    def handle_data(self, data: str) -> None:
        if self.body_found:
            self.text.write(data)

    def get_data(self) -> str:
        return self.text.getvalue().strip()


def strip_tags(html: str):
    s = MLStripper()
    s.feed(html)
    return s.get_data()


if __name__ == "__main__":
    """
    This program converts every *.html file that is considered for
    an HTML email to have their *.txt counterparts with the same name,
    just with a different file extension.
    
    If the *.html file already have its' counterparts, we will skip that one.
    """
    files_collection: dict[str, bool] = {}
    for file in listdir("."):
        if not isfile(file):
            continue

        if file.endswith(".html") or file.endswith(".txt"):
            files_collection[file] = False

    for filename, state in files_collection.items():
        if filename.endswith(".html"):
            # Does it have *.txt counterparts?
            raw_filename = filename[:-5]
            if files_collection.get(raw_filename + ".txt") is not None:
                # We have it, skipping
                print(f"{filename} already got text counterpart")
                continue

            with open(filename, "r", encoding='utf-8') as html_file:
                file_content = html_file.read()
                removed_tags = strip_tags(file_content)

                with open(raw_filename + ".txt", "w", encoding='utf-8') as text_file:
                    text_file.write(removed_tags)

                print(f"Removed HTML tags for {raw_filename}")
