#!/bin/bash

while read e; do
  echo $e
  curl https://fonts.gstatic.com/s/e/notoemoji/latest/$e/lottie.json --output ui/public/emojis/lottie/$e.json
done < emojis.txt

# while read e; do
#   echo ffmpeg -i ui/public/emojis/512/$e.webp -vf "scale=64:64" ui/public/emojis/64/$e.webp
#   ffmpeg -i ui/public/emojis/512/$e.webp -vf "scale=64:64" ui/public/emojis/64/$e.webp
# #   curl https://fonts.gstatic.com/s/e/notoemoji/latest/$e/512.webp --output ui/public/emojis/512/$e.webp
# done < emojis.txt

# magick ui/public/emojis/512/1f60a.webp -coalesce -resize 128x -layers optimize ui/public/emojis/64/1f60a.webp
# magick  ui/public/emojis/512/1f60a.webp -coalesce -resize 64x -quality 50 -strip -define webp:method=6 -layers optimize ui/public/emojis/64/1f60a.webp