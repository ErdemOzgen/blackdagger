# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'Blackdagger'
copyright = '2024, Blackdagger Developers'
author = 'Blackdagger Developers'
release = '1.03'

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = []

templates_path = ['_templates']
exclude_patterns = []


# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output


html_theme = 'sphinx_rtd_theme'
html_static_path = ['_static']
html_css_files = ['css/custom.css']
html_logo = '_static/logo.png'
html_theme_options = {
    'style_external_links': True,
    'style_nav_header_background': '#333333',  # Header background
}

# -- Options for Localization ------------------------------------------------
locale_dirs = ['locale/']  
gettext_compact = False  


# layout
templates_path = ['_templates']