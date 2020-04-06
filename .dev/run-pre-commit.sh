#!/usr/bin/env sh

if [ -n "$SKIP_PRE_COMMIT" ]; then
  echo "✔️✔️✔️ Skipping pre-commit because env var SKIP_PRE_COMMIT exists and not-empty ✔️✔️✔️"
  exit 0
fi

CHANGED_GO_FILES=$(git diff HEAD --name-only | grep -P '.+\.go$')

if [ -z "$CHANGED_GO_FILES" ]; then

  echo "✌️ No golang files changed ✌️"

else

  echo "👉 Check Dep 👈"
  if ! make check-dep; then
    echo "⛔ Dep's Gopkg inconsistent or dep version not above 0.5! ⛔"
    exit 1
  fi
  echo "✔️ Dep OK ✔️"

  echo "👉 Linting 👈"
  if ! make lint; then
    echo "⛔ Code unclean, linting failed ⛔"
    exit 1
  fi
  echo "✔️ Lint OK ✔️"

  echo "👉 Testing 👈"
  if ! make test; then
    echo "⛔ Test failed, code not robust! ⛔"
    exit 1
  fi
  echo "✔️ Test OK ✔️"

  echo "👉 Check Buildable 👈"
  if ! make check-buildable; then
    echo "⛔ Build failed! ⛔"
    exit 1
  fi
  echo "✔️ Build OK ✔️"

fi

echo "✔️✔️✔️ Pre-Commit OK ✔️✔️✔️"
exit 0