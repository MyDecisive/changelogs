set -e 
tracked=$(yq ".version" tracked.yaml)
if [  "$1" = "`echo -e "$1\n$tracked" | sort -V | head -n1`" ]; then
  echo "Already generated"
else
  export version=$1
  composite --config https://raw.githubusercontent.com/DecisiveAI/changelogs/refs/heads/ENG-972/scripts/composite/cliff.toml --path $GITHUB_WORKSPACE/CHANGELOG.md --version $version
  yq e -i '.version= env(version)' tracked.yaml
  echo "GENERATED=1" >> "$GITHUB_OUTPUT"
fi
