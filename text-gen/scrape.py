from requests import get
from bs4 import BeautifulSoup as bs
import re
import os
import pdb


url_season = "https://genius.com/albums/How-i-met-your-mother/Season-1"
response_season = get(url_season)

url_whole = "https://genius.com/artists/How-i-met-your-mother"
response_whole = get(url_whole)

#Try to get the whole show urls now
soup_whole = bs(response_whole.text,"html.parser")
url_whole = soup_whole.find_all('div', class_ = "modal_window-content modal_window-content--narrow_width modal_window-content--white_background")
lnk = str(soup_whole).strip(",")
links = [i for i in lnk if "href" in i]
#links[7] contains the urls
#https://genius.com/albums/How-i-met-your-mother/Season-(number)
#except for S08 for some reason
whole_links = ["https://genius.com/albums/How-i-met-your-mother/Season-{}".format(i) for i in range(1, 9)]

all_links = []

for link in whole_links:
    response_season = get(link)
    soup_season = bs(response_season.text, 'html.parser')
    link_season = soup_season.find_all('div', class_ = 'column_layout-column_span column_layout-column_span--primary')
    #len(ly_season = 47) from 3, every other item
    ep_content = []
    for x in link_season: 
        ep_content.append(x)
    #something[3] is the pilot stuff, look at 
    ep_content = str(ep_content).split("\n")
    link_location = [i for i in ep_content if "href" in i]

    # All season 1 links
    for class_ in link_location:
        all_links.append(re.search(r"(?<=href=\").*?(?=\">)", class_).group(0))

#[0] #all episodes heregit 
print("links generated")
barney_lines = []
for url in all_links:
    response = get(url)

    html_soup = bs(response.text, 'html.parser')
    # maybe .text
    lyrics_cont = html_soup.find_all('div', class_ = 'song_body-lyrics')
    if not lyrics_cont:
        continue
    print("Found some lyrics")
    lyrics = str(lyrics_cont[0]).split("\n")
    #barney = [x for x in lyrics if "Barney:" in x]
    #if barney:
    #    print("With barney lines")
   # pdb.set_trace()
    barney_lines.append(lyrics)

html_items = r"<.*>"
click_item = r"ng-click.*\)\""
rando = r"pending.*\">"

clean = re.compile(r"{}|{}|{}".format(html_items, click_item, rando))
print("Cleaning output")
for i, line in enumerate(barney_lines):
    
    #line = line.remove("<br/>").remove("</div>")
    barney_lines[i] = [re.sub(clean, '', str(one_line)) for one_line in line if re.sub(clean, '', str(one_line))]
    #barney_lines[i] = [x of x in ]
print(len(barney_lines))

#robots vs wrestlers, im glad u asked ted
with open("lines.txt", "w") as f:
    for batch in barney_lines:
        for line in batch:
            f.write(line + "\n")

print("Done!!")
