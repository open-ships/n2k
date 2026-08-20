#!/usr/bin/env bash

set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

version=$(< VERSION)
if [[ ! "$version" =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
    echo "VERSION must contain an exact semantic version without a leading v: $version" >&2
    exit 1
fi

first_heading=$(grep -m 1 '^### ' CHANGELOG.md || true)
heading_pattern="^### v${version} — [0-9]{4}-[0-9]{2}-[0-9]{2} — .+$"
if [[ ! "$first_heading" =~ $heading_pattern ]]; then
    echo "CHANGELOG.md first release heading must match:" >&2
    echo "### v$version — YYYY-MM-DD — <summary>" >&2
    echo "found: ${first_heading:-<none>}" >&2
    exit 1
fi

target_tag="v$version"
head_commit=$(git rev-parse HEAD)
if git show-ref --verify --quiet "refs/tags/$target_tag"; then
    if [[ $(git cat-file -t "refs/tags/$target_tag") != tag ]]; then
        echo "$target_tag exists but is not an annotated tag" >&2
        exit 1
    fi
    tag_commit=$(git rev-list -n 1 "$target_tag")
    if [[ "$tag_commit" != "$head_commit" ]]; then
        echo "$target_tag already points at $tag_commit; bump VERSION for $head_commit" >&2
        exit 1
    fi
    echo "release metadata matches existing $target_tag at $head_commit"
    exit 0
fi

latest_tag=""
while IFS= read -r candidate; do
    if [[ "$candidate" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
        latest_tag=$candidate
        break
    fi
done < <(git tag --list 'v*' --sort=-version:refname)
if [[ -n "$latest_tag" ]]; then
    latest_version=${latest_tag#v}
    if ! awk -v proposed="$version" -v latest="$latest_version" 'BEGIN {
        split(proposed, p, ".")
        split(latest, l, ".")
        for (i = 1; i <= 3; i++) {
            if ((p[i] + 0) > (l[i] + 0)) exit 0
            if ((p[i] + 0) < (l[i] + 0)) exit 1
        }
        exit 1
    }'; then
        echo "VERSION $version must be newer than latest published tag $latest_tag" >&2
        exit 1
    fi
fi

echo "release metadata is ready for $target_tag at $head_commit"
