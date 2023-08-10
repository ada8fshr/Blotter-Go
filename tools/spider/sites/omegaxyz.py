import requests
from bs4 import BeautifulSoup
if __name__ == "__main__":
    from utils import *
else:
    from .utils import *


class OmegaXYZ(Site):
    def __init__(self):
        super(Site, self)

    def matcher(self, url: str):
        return url == "https://www.omegaxyz.com/"

    def solver(self, url: str):
        res = get("%s/archive" % url.strip("/"))
        soup = BeautifulSoup(res, features="lxml")
        posts = []
        for item in soup.select("h3.rpwe-title"):
            link = item.select_one("a")
            posts.append(Post(
                link.get_text(),
                link.get("href"),
                parseToUnix(item.parent.select_one("time").get("datetime")),
            ))
        return posts


if __name__ == '__main__':
    t = OmegaXYZ()
    print(t.matcher("https://www.omegaxyz.com/"))
    print(t.solver("https://www.omegaxyz.com/"))
