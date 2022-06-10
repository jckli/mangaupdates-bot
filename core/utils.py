import re
from bs4 import BeautifulSoup as bs

class Util:
    def __init__(self):
        pass

    def format_mu_description(self, description):
        description = re.sub(r"\b<BR>\b", " ", description)
        soup = bs(description, "html.parser")
        soup = soup.text
        soup = " ".join(soup.split())
        return soup

    def format_group_name(self, group_name):
        group_name = re.sub(r"&#039;", "'", group_name)
        return group_name