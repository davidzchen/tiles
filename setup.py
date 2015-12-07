# Copyright 2014 David Z. Chen
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not
# use this file except in compliance with the License. You may obtain a copy of
# the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations under
# the License.

from setuptools import setup

setup(
    name = "tiles",
    version = "1.0.0",
    url = "https://github.com/davidzchen/tiles",
    author = "David Z. Chen",
    author_email = "david@davidzchen.com",
    description = "Easy tmux management",
    long_description = open("README.md").read(),
    scripts = ["tiles"],
)
