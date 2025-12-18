set -e 
tracked=$(yq ".version" tracked.yaml)
if [  "$1" = "`echo -e "$1\n$tracked" | sort -V | head -n1`" ]; then
  echo "Already generated"
else
  export version=$1
  composite --path $GITHUB_WORKSPACE/CHANGELOG.md --version $version
  yq e -i '.version= env(version)' tracked.yaml
  echo "GENERATED=1" >> "$GITHUB_OUTPUT"
fi
