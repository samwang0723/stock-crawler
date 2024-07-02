#!/bin/bash

# Function to get the latest tag
get_latest_tag() {
  git describe --tags `git rev-list --tags --max-count=1`
}

# Function to increment the version number
increment_version() {
  local version=$1
  local major=$(echo $version | cut -d. -f1 | cut -c2-)
  local minor=$(echo $version | cut -d. -f2)
  local patch=$(echo $version | cut -d. -f3)

  patch=$((patch + 1))

  echo "v$major.$minor.$patch"
}

# Get the latest tag
latest_tag=$(get_latest_tag)

# If no tags are found, start with v1.0.0
if [ -z "$latest_tag" ]; then
  new_tag="v1.0.0"
else
  new_tag=$(increment_version $latest_tag)
fi

# Replace old version tag with new one in specified files
files_to_update=("Makefile" "deployments/helm/stock-crawler/deployment.yaml" "configs/config.dev.yaml" "configs/config.prod.yaml" "build/docker/app/docker-compose.yml")

# Include file version changes in commit
for file in "${files_to_update[@]}"; do
  if [ -f "$file" ]; then
    sed -i '' "s/$latest_tag/$new_tag/g" "$file"
  fi
  git add $file
done

git commit -m "Bump version to $new_tag"

# Create a new tag
git tag $new_tag

# Generate changelogs
make changelog-gen
make changelog-commit

# Push the new tag to the remote repository
git push origin $new_tag
git push origin master

echo "New tag $new_tag created and changelogs generated."

