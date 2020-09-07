from setuptools import find_namespace_packages
from setuptools import setup

REQUIRED_PACKAGES = [
    "grpcio", "protobuf", "google-api-core",
]

setup(
    name='barapis',
    version='1.0',
    author='Hayo van Loon',
    author_email='hayovanloon@gmail.com',
    install_requires=REQUIRED_PACKAGES,
    package_dir={'': 'var'},
    packages=find_namespace_packages(where='var'),
    include_package_data=True,
    description='Protocol Buffer specifications for the Bar application',
    requires=[],
)
