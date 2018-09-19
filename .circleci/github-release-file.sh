#!/bin/sh

set -eu
set -x

repo="$1"; shift
tag="$1"; shift
filename="$1"; shift
content_type="$1"; shift

creds="${GITHUB_USER}:${GITHUB_ACCESS_TOKEN}"

test -z "$repo" -o -z "$tag" -o -z "$filename" -o -z "$content_type" -o -z "$creds" && exit 1

echo 'Creating release ...'
response=$(curl -s -S \
                -u "$creds" \
                -X POST \
                -T- \
                "https://api.github.com/repos/${repo}/releases" <<EOF)
{
  "tag_name": "$tag",
  "draft": true
}
EOF

# remove hypermedia garbage
eval "$(jq -r '"upload_url=\( .upload_url | sub("{.*}$";"") | @sh )"' <<EOF)"
$response
EOF

echo 'Uploading file ...'
response=$(curl -s -S \
                -u "$creds" \
                -X POST \
                -H "Content-Type: $content_type" \
                -T "$filename" \
                "${upload_url}?name=${filename}")

eval "$(jq -r '"state=\(.state | @sh) size=\(.size | @sh) url=\(.browser_download_url | @sh)"' <<EOF)"
$response
EOF

if [ "$state" = "uploaded" ]; then
    echo "Done ($size bytes)."
    echo "'$filename' is available as <$url>."
else
    echo "Error: state=$state"
    exit 1
fi
