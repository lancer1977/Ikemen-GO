import pathlib
import sys
s = pathlib.Path('.github/workflows/local-deploy.yml').read_text()
checks = [
    ('workflow_dispatch', 'workflow_dispatch:' in s),
    ('linux labels', 'runs-on: [self-hosted, linux, ikemen-linux]' in s),
    ('windows labels', 'runs-on: [self-hosted, windows, ikemen-windows]' in s),
    ('msys2 shell', 'shell: msys2 {0}' in s),
    ('windows wrapper', './scripts/deploy-local-windows.sh' in s),
    ('linux wrapper', './scripts/deploy-local-linux.sh' in s),
]
for name, ok in checks:
    print(f'{name}: {ok}')
sys.exit(0 if all(ok for _, ok in checks) else 1)
