---
title: Text Search
parent: Tutorials
layout: default
nav_order: 1
---

# Text Search

Text search has been designed to "just work", but understanding the available syntax and features can help you get the most out of it.

## Quotes

Unquoted search terms will match any record containing all the terms, in any order. Terms (and parts of terms) can be enclosed in quotes to specify that they must appear in the order specified. For example, searching for `"fat cats and sad rats"` will return results that contain this exact phrase, whereas searching for `fat cats and sad rats` might return results that contain the phrase "rats and fat, sad cats".

A common pattern would be searching for a quoted movie title along with an unquoted year, e.g. `"steamboat willie" 1928`.

## Followed By Operator

The `.` character is an alternative way of specifying word order. For example, searching for `apple . orange` will return results that contain `apple` followed by `orange`.

## OR Operator

The `|` character is an OR operator. Search will return results that contain either term on either side of the `|`. For example, searching for `apple | orange` will return results that contain either `apple` or `orange`.

## Negation Operator

The `!` character is a negation operator. Search will return results that do not contain the term following the `!`. For example, searching for `orange !apple` will return results that contain `orange` _and_ do not contain `apple`.

## Wildcard Suffix

The `*` character can be used as a wildcard suffix. Search will return results that start with the term preceding the `*`. For example, searching for `appl*` will return results that start with `appl`, such as `apple`, `application`, `appliance`, etc.

Note that wildcards can only be used as a suffix, not a prefix or infix.

## Parentheses

Parentheses can be used to control operator precedence. For example, searching for `"banana split" | (apple !toffee)` will return results that either contain the phrase "banana split", or contain the word "apple", but not the word "toffee".

## Normalisation

Search text is case-insensitive. Any punctuation characters not mentioned above are treated as whitespace. The search index along with search terms are normalised to ascii characters (meaning searches for `cafe` will match records containing "café", and vice versa). This applies to all languages, for example simplified and traditional Chinese characters can be considered interchangeable for search.

## Syntax validation

The search query parser is forgiving: any string is valid text search input, and the parser will make a best effort to apply the syntax described above to any input string. For example if you search for `"Steamboat Willie` (forgetting to close the quotes), the search query parser will still consider this to be a quoted phrase.
