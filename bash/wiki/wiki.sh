#!/bin/bash

#### Constants

WIKI_URL="https://en.wikipedia.org/wiki/"

#### Main

# Gather command line arguments

if [ "$#" = "0" ]; then
  echo -e "Must provide article name, e.g. ./wiki.sh walrus"
  exit 1
fi
article=$(echo "$1" | awk '{for (i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)} 1' | tr ' ' '_')
subsection=$(echo "$2" | awk '{for (i=1;i<=1;i++) $i=toupper(substr($i,1,1)) substr($i,2)} 1' | tr ' ' '_')

echo -e "Article name: $article";
echo -e "Subsection: $subsection";

# Grab data

data_file="resources/$article.html"
use_local=
if [ -f "resources/$article.html" ]; then
    echo "using local $data_file copy."
    use_local=true
else
    echo "fetching from $WIKI_URL$article"
#    curl -o "$data_file" "$WIKI_URL$article"
fi

echo $use_local

article_html=$(if [ "$use_local" ]; then cat "$data_file"; else curl -s "$WIKI_URL$article"; fi)

if [[ "$subsection" = "" ]]
then
  # First legit sentence of the article lies within a div element with a "mw-parser-output" class and is usually inside of the first non-empty <p> tag
  echo "$article_html" | grep -A 1000 '<div class="mw-parser-output"' | grep -i "$article" | grep "<p>" | head -n 1 | cut -d "." -f1 | sed -e 's/<[^>]*>//g'

  echo "$article_html" | grep "toclevel-1" | sed -e 's/<[^>]*>//g'
else
  echo "$article_html" | grep -A 5 "id=\"$subsection"\" | grep "<p>" | head -n 1 | cut -d "." -f1 | sed -e 's/<[^>]*>//g'

fi

