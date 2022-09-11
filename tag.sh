git push
git tag  $(date +"%Y-%m-%d.%Hh%Mm")-$(git rev-parse --short HEAD)
git push --tags