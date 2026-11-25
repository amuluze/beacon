#!/usr/bin/env python3
"""Fail closed on low-quality or cross-module generated signal entries."""

import glob
import re
import sys


SECTION_NAMES = ('接口与符号信号', '配置与入口信号', '数据模型与状态')
CHILD_MODULES = {
    '.docs/modules/common.md': ('common/docker/', 'common/logger/', 'common/timex/'),
}


def section_lines(body, name):
    match = re.search(r'^## ' + re.escape(name) + r'\n(.*?)(?=^## |\Z)', body, re.M | re.S)
    if not match:
        return []
    return [line for line in match.group(1).split('\n') if line.startswith('- ')]


def signal_keys(items):
    keys = set()
    for line in items:
        match = re.search(r'`([^`]+)`', line)
        if match:
            keys.add(match.group(1))
    return keys


def main():
    issues = []
    docs = sorted(glob.glob('.docs/modules/*.md'))
    for doc in docs:
        body = open(doc, encoding='utf-8').read()
        sections = {name: section_lines(body, name) for name in SECTION_NAMES}
        interface = sections['接口与符号信号']
        config = sections['配置与入口信号']
        interface_keys = signal_keys(interface)
        config_keys = signal_keys(config)
        if interface_keys and config_keys:
            overlap = interface_keys & config_keys
            if len(overlap) / max(len(interface_keys), len(config_keys)) > 0.5:
                issues.append(f'{doc}: duplicate interface/config signals: {len(overlap)}')
        for line in interface + config:
            match = re.search(r'[A-Za-z0-9_./-]+(?:_test\.go|\.test\.ts|\.spec\.ts)', line)
            if match:
                issues.append(f'{doc}: Test file in entry signal: {match.group(0)}')
        for prefix in CHILD_MODULES.get(doc, ()):
            for name, lines in sections.items():
                for line in lines:
                    if any(token.startswith(prefix) for token in re.findall(r'`([^`]+)`', line)):
                        issues.append(f'{doc}: Child module signal in {name}: {prefix}')
    if issues:
        print('\n'.join(issues))
        return 1
    print(f'{len(docs)} module documents passed signal ownership and test-file checks')
    return 0


if __name__ == '__main__':
    sys.exit(main())
