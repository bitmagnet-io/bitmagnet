---
title: Classifier
parent: Guides
layout: default
nav_order: 4
redirect_from:
  - /tutorials/classifier.html
---

# Classifier

{: .note-title }

> tl;dr:
>
> The classifier can be configured and customized to do things like:
>
> - automatically delete torrents you don't want in your index
> - add custom tags to torrents you're interested in
> - customize the keywords and file extensions used for determining a torrent's content type
> - specify completely custom logic to classify and perform other actions on torrents
>
> Skip to [practical use cases and examples](#practical-use-cases-and-examples)

## Background

After a torrent is crawled or imported, some further processing must be done to gather metadata, have a guess at the torrent's contents and finally index it in the database, allowing it to be searched and displayed in the UI/API.

**bitmagnet**'s classifier is powered by a [Domain Specific Language](https://en.wikipedia.org/wiki/Domain-specific_language). The aim of this is to provide a high level of customisability, along with transparency into the classification process which will hopefully aid collaboration on improvements to the core classifier logic.

The classifier is declared in YAML format. The application includes a [core classifier](https://github.com/bitmagnet-io/bitmagnet/blob/main/internal/classifier/classifier.core.yml) that can be configured, extended or completely replaced with a custom classifier. This page documents the required format.

## Source precedence

**bitmagnet** will attempt to load classifier source code from all the following locations. Any discovered classifier source will be merged with other sources in the following order of precedence:

- [the core classifier](https://github.com/bitmagnet-io/bitmagnet/blob/main/internal/classifier/classifier.core.yml)
- `classifier.yml` in the [XDG-compliant](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) config location for the current user (for example on MacOS this is `~/Library/Application Support/bitmagnet/classifier.yml`)
- `classifier.yml` in the current working directory
- [Classifier configuration](#configuration)

Note that multiple sources will be merged, not replaced. For example, keywords added to the classifier configuration will be merged with the core keywords.

The merged classifier source can be viewed with the CLI command `bitmagnet classifier show`.

{% include callout_cli.md %}

## Schema

A [JSON schema for the classifier](https://bitmagnet.io/schemas/classifier-0.1.json) is available; some editors and IDEs will be able to validate the structure of your classifier document by specifying the `$schema` attribute:

```yaml
$schema: https://bitmagnet.io/schemas/classifier-0.1.json
```

The classifier schema can also be viewed by running the cli command `bitmagnet classifier schema`.

{% include callout_cli.md %}

The classifier declaration comprises the following components:

## Workflows

A workflow is a list of [actions](#actions) that will be executed on all torrents when they are classified. When no custom configuration is provided, the `default` workflow will be run. To use a different workflow instead, specify the `classifier.workflow` configuration option with the name of your custom workflow.

## Actions

An action is a piece of [workflow](#workflows) to be executed. All actions either return an updated classification result or an error.

For example, the following action will set the content type of the current torrent to `audiobook`:

```yaml
set_content_type: audiobook
```

The following action will return an `unmatched` error:

```yaml
unmatched
```

And the following action will delete the current torrent being classified (returning a `delete` error):

```yaml
delete
```

These actions aren't much use on their own - we'd want to check some conditions are satisfied before setting a content type or deleting a torrent, and for this we'd use the `if_else` action. For example, the following action will set the content type to `audiobook` if the torrent name contains audiobook-related keywords, and will otherwise return an `unmatched` error:

```yaml
if_else:
  condition: "torrent.baseName.matches(keywords.audiobook)"
  if_action:
    set_content_type: audiobook
  else_action: unmatched
```

The following action will delete a torrent if its name matches the list of`banned` keywords:

```yaml
if_else:
  condition: "torrent.baseName.matches(keywords.banned)"
  if_action: delete
```

Actions may return the following types of error:

- An `unmatched` error indicates that the current action did not match for the current torrent
- A `delete` error indicates that the torrent should be deleted
- An unhandled error may occur, for example if the TMDB API was unreachable

Whenever an error is returned, the current classification will be terminated.

Note that a workflow should never return an `unmatched` error. We expect to iterate through a series of checks corresponding to each content type. If the current torrent does not match the content type being checked, we'll proceed to the next check until we find a match; if no match can be found, the content type will be `unknown`. To facilitate this, we can use the `find_match` action.

The `find_match` action is a bit like a try/catch block in some programming languages; it will try to match a particular content type, and if an `unmatched` error is returned, it will catch the `unmatched` error proceed to the next check. For example, the following action will attempt to classify a torrent as an `audiobook`, and then as an `ebook`. If both checks fail, the content type will be `unknown`:

```yaml
find_match:
  # match audiobooks:
  - if_else:
      condition: "torrent.baseName.matches(keywords.audiobook)"
      if_action:
        set_content_type: audiobook
      else_action: unmatched
  # match ebooks:
  - if_else:
      condition: "torrent.files.map(f, f.extension in extensions.ebook ? f.size : - f.size).sum() > 0"
      if_action:
        set_content_type: ebook
      else_action: unmatched
```

For a full list of available actions, please refer to [the JSON schema](https://bitmagnet.io/schemas/classifier-0.1.json).

## Conditions

Conditions are used in conjunction with the `if_else` [action](#actions), in order to execute an action if a particular condition is satisfied.

The conditions in the examples above use [CEL (Common Expression Language) expressions](https://cel.dev/).

### The CEL environment

CEL is already a [well-documented](https://github.com/google/cel-spec/blob/master/doc/intro.md) language, so this page won't go into detail about the CEL syntax. In the context of the **bitmagnet** classifier, the CEL environment exposes a number of variables:

- `torrent`: The current torrent being classified (protobuf type: `bitmagnet.Torrent`)
- `result`: The current classification result (protobuf type: `bitmagnet.Classification`)
- `keywords`: A map of strings to regular expressions, representing named lists of [keywords](#keywords)
- `extensions`: A map of strings to string lists, representing named lists of [extensions](#extensions)
- `contentType`: A map of strings to enum values representing content types (e.g. `contentType.movie`, `contentType.music`)
- `fileType`: A map of strings to enum values representing file types (e.g. `fileType.video`, `fileType.audio`)
- `flags`: A map of strings to the configured values of [flags](#flags)
- `kb`, `mb`, `gb`: Variables defined for convenience, equal to the number of bytes in a kilobyte, megabyte and gigabyte respectively

For more details on the protocol buffer types, please refer to [the protobuf schema](https://github.com/bitmagnet-io/bitmagnet/blob/main/internal/protobuf/bitmagnet.proto).

### Boolean logic (`or`, `and` & `not`)

In addition to CEL expressions, conditions may be declared using the boolean logic operators `or`, `and` and `not`. For example, the following condition evaluates to true, if either the torrent consists mostly of file extensions very commonly used for music (e.g. `flac`), OR if the torrent both has a name that includes music-related keywords, and consists mostly of audio files:

```yaml
or:
  - "torrent.files.map(f, f.extension in extensions.music ? f.size : - f.size).sum() > 0"
  - and:
      - "torrent.baseName.matches(keywords.music)"
      - "torrent.files.map(f, f.fileType == fileType.audio ? f.size : - f.size).sum() > 0"
```

Note that we could also have specified the above condition using just one CEL expression, but breaking up complex conditions like this is more readable.

## Keywords

The classifier includes lists of keywords associated with different types of torrents. These aim to provide a simpler alternative to regular expressions, and the classifier will compile all keyword lists to regular expressions that can be used within CEL expressions. In order for a keyword to match, it must appear as an isolated token in the test string - that is, it must be either at the beginning or preceded by a non-word character, and either at the end or followed by a non-word character.

Reserved characters in the syntax are:

- parentheses `(` and `)` enclose a group
- `|` is an OR operator
- `*` is a wildcard operator
- `?` makes the previous character or group optional
- `+` specifies one or more of the previous character
- `#` specifies any number
- ` ` specifies any non-word or non-number character

For example, to define some music- and audiobook-related keywords:

```yaml
keywords:
  music: # define music-related keywords
    - music # all letters are case-insensitive, and must be defined in lowercase unless escaped
    - discography
    - album
    - \V.?\A # escaped letters are case-sensitive; matches "VA", "V.A" and "V.A.", but not "va"
    - various artists # matches "various artists" and "Various.Artists"
  audiobook: # define audiobook-related keywords
    - (audio)?books?
    - (un)?abridged
    - narrated
    - novels?
    - (auto)?biograph(y|ies) # matches "biography", "autobiographies" etc.
```

{: .note }

> If you'd rather use plain old regular expressions, the CEL syntax supports that too, for example `torrent.baseName.matches("^myregex$")`.

## Extensions

The classifier includes lists of file extensions associated with different types of content. For example, to identify torrents of type `comic` by their file extensions, the extensions are first declared:

```yaml
extensions:
  comic:
    - cb7
    - cba
    - cbr
    - cbt
    - cbz
```

The extensions can now be used as part of a condition within an `if_else` action:

```yaml
if_else:
  condition: "torrent.files.map(f, f.extension in extensions.comic ? f.size : - f.size).sum() > 0"
  if_action:
    set_content_type: comic
  else_action: unmatched
```

## Flags

Flags can be used to configure workflows. In order to use a flag in a workflow, it must first be defined. For example, the core classifier defines the following flags that are used in the `default` workflow:

```yaml
flag_definitions:
  tmdb_enabled: bool
  delete_content_types: content_type_list
  delete_xxx: bool
```

These flags can be referenced within CEL expressions, for example to delete adult content if the `delete_xxx` flag is set to `true`:

```yaml
if_else:
  condition: "flags.delete_xxx && result.contentType == contentType.xxx"
  if_action: delete
```

## Configuration

The classifier can be customized by providing a `classifier.yml` file in a supported location [as described above](#source-precedence). If you only want to make some minor modifications, it may be convenient to specify these [using the main application configuration](/setup/configuration.html) instead, by providing values in either `config.yml` or as environment variables. The application configuration exposes some but not all properties of the classifier.

For example, in your `config.yml` you could specify:

```yaml
classifier:
  # specify a custom workflow to be used:
  workflow: custom
  # add to the core list of music keywords:
  keywords:
    music:
      - my-custom-music-keyword
  # add a file extension to the list of audiobook-related extensions:
  extensions:
    audiobook:
      - abc
  # auto-delete all comics
  flags:
    delete_content_types:
      - comics
```

Or as environment variables you could specify:

```sh
TMDB_ENABLED=false \ # disable the TMDB API integration
  CLASSIFIER_WORKFLOW=custom \ # specify a custom workflow to be used
  CLASSIFIER_DELETE_XXX=true \ # auto-delete all adult content
  bitmagnet worker run --all
```

## Validation

The classifier source is compiled on initial load, and all structural and syntax errors should be caught at compile time. If there are errors in your classifier source, **bitmagnet** should exit with an error message indicating the location of the problem.

## Testing on individual torrents

You can test the classifier on an individual torrent or torrents using the `bitmagnet process` CLI command:

```sh
bitmagnet process --infoHash=aaaaaaaaaaaaaaaaaaaa --infoHash=bbbbbbbbbbbbbbbbbbbb
```

{% include callout_cli.md %}

## Reclassify all torrents

Read how to [reclassify all torrents](/guides/reprocess-reclassify.html).

## Practical use cases and examples

### Auto-delete specific content types

The default workflow provides a flag that allows for automatically deleting specific content types. For example, to delete all `comic`, `software` and `xxx` torrents:

```yaml
flags:
  delete_content_types:
    - comic
    - software
    - xxx
```

Auto-deleting adult content has been one of the most requested features. For convenience, this is exposed as the configuration option `classifier.delete_xxx`, and can be specified with the environment variable `CLASSIFIER_DELETE_XXX=true`.

### Auto-delete torrents containing specific keywords

Any torrents containing keywords in the `banned` list will be automatically deleted. This is primarily used for deleting <abbr title="Child Sexual Abuse Material">CSAM</abbr> content, but the list can be extended to auto-delete any other keywords:

```yaml
keywords:
  banned:
    - my-hated-keyword
```

### Disable the TMDB API integration

The `tmdb_enabled` flag can be used to disable the TMDB API integration:

```yaml
flags:
  tmdb_enabled: false
```

For convenience, this is also exposed as the configuration option `tmdb.enabled`, and can be specified with the environment variable `TMDB_ENABLED=false`.

The `apis_enabled` flag has the same effect, disabling TMDB and any future API integrations:

```yaml
flags:
  apis_enabled: false
```

API integrations can also be disabled for individual classifier runs, without disabling them globally, by passing the `--apisDisabled` flag to [the reprocess command](/guides/reprocess-reclassify.html).

### Extend the default workflow with custom logic

Custom workflows can be added in the `workflows` section of the classifier document. It is possible to extend the default workflow by using the `run_workflow` action within your custom workflow, for example:

```yaml
workflows:
  custom:
    - <my custom action to be executed before the default workflow>
    - run_workflow: default
    - <my custom action to be executed after the default workflow>
```

A concrete example of this is adding tags to torrents based on custom criteria.

### Use tags to create custom torrent categories

Is there a category of torrent you're interested in that isn't captured by one of the core content types? Torrent tags are intended to capture custom categories and content types.

Let's imagine you'd like to surface torrents containing interesting documents. The interesting documents have specific file extensions, and their filenames contain specific keywords. Let's create a custom action to tag torrents containing interesting documents:

```yaml
# define file extensions for the documents we're interested in:
extensions:
  interesting_documents:
    - doc
    - docx
    - pdf
# define keywords that must be present in the filenames of the interesting documents:
keywords:
  interesting_documents:
    - interesting
    - fascinating
# extend the default workflow with a custom workflow to tag torrents containing interesting documents:
workflows:
  custom:
    # first run the default workflow:
    - run_workflow: default
    # then add the tag to any torrents containing interesting documents:
    - if_else:
        condition: "torrent.files.filter(f, f.extension in extensions.interesting_documents && f.basePath.matches(keywords.interesting_documents)).size() > 0"
        if_action:
          add_tag: interesting-documents
```

To specify that the custom workflow should be used, remember to specify the `classifier.workflow` configuration option, e.g. `CLASSIFIER_WORKFLOW=custom bitmagnet worker run --all`.
