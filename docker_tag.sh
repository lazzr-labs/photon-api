# !/bin/bash

docker_tag="ute1787929769"
epoch_time="ute$(date +%s)"

files=(
  "docker_buildx.sh"
  "docker_tag.sh"
  "docker-api.yaml"
  "helm-v1/values.yaml"
)

for file in "${files[@]}"; do
  sed -i '' -e "s/$docker_tag/$epoch_time/g" "$file"
done
