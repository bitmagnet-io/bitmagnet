package gen

import (
	"github.com/iancoleman/strcase"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"path"
	"runtime"
)

func BuildGenerator(db *gorm.DB) *gen.Generator {
	_, filename, _, _ := runtime.Caller(0)
	internal := path.Dir(path.Dir(path.Dir(filename)))

	g := gen.NewGenerator(gen.Config{
		OutPath:       internal + "/database/dao",
		ModelPkgPath:  internal + "/model",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		FieldSignable: true,
	})

	g.UseDB(db)

	g.WithDataTypeMap(map[string]func(gorm.ColumnType) (dataType string){
		"text": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "NullString"
			}
			return "string"
		},
		"bytea": func(columnType gorm.ColumnType) (dataType string) {
			return "[]byte"
		},
	})

	infoHashType := gen.FieldType("info_hash", "protocol.ID")
	infoHashReadOnly := readAndCreateField("info_hash")
	createdAtReadOnly := readAndCreateField("created_at")

	g.WithJSONTagNameStrategy(strcase.ToLowerCamel)

	torrentSources := g.GenerateModel(
		"torrent_sources",
		readAndCreateField("key"),
		createdAtReadOnly,
	)
	torrentFiles := g.GenerateModel(
		"torrent_files",
		infoHashType,
		infoHashReadOnly,
		gen.FieldType("size", "uint"),
		gen.FieldType("index", "uint"),
		gen.FieldGORMTag("index", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		gen.FieldGORMTag("path", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		gen.FieldGenType("extension", "String"),
		gen.FieldGORMTag("extension", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "false")
			return tag
		}),
		createdAtReadOnly,
	)
	torrentPieces := g.GenerateModelAs(
		"torrent_pieces",
		"TorrentPieces",
		infoHashType,
		infoHashReadOnly,
		gen.FieldJSONTag("pieces", "-"),
		createdAtReadOnly,
	)
	torrentsTorrentSources := g.GenerateModel(
		"torrents_torrent_sources",
		readAndCreateField("source"),
		infoHashType,
		infoHashReadOnly,
		gen.FieldType("seeders", "NullUint"),
		gen.FieldType("leechers", "NullUint"),
		gen.FieldType("published_at", "sql.NullTime"),
		gen.FieldRelate(
			field.HasOne,
			"TorrentSource",
			torrentSources,
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": []string{"Source"},
				},
			},
		),
		createdAtReadOnly,
	)
	torrentTags := g.GenerateModel(
		"torrent_tags",
		infoHashType,
		infoHashReadOnly,
		gen.FieldGenType("name", "String"),
		gen.FieldGORMTag("name", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		createdAtReadOnly,
	)
	torrents := g.GenerateModel(
		"torrents",
		gen.FieldRelate(
			field.HasOne,
			"Hint",
			g.GenerateModel(
				"torrent_hints",
				infoHashType,
			),
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": []string{"InfoHash"},
				},
			},
		),
		gen.FieldRelate(
			field.HasMany,
			"Contents",
			g.GenerateModel(
				"torrent_contents",
				infoHashType,
			),
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": []string{"InfoHash"},
				},
			},
		),
		gen.FieldRelate(
			field.HasMany,
			"Sources",
			torrentsTorrentSources,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag: field.GormTag{
					"foreignKey": []string{"InfoHash"},
				},
			},
		),
		gen.FieldRelate(
			field.HasMany,
			"Files",
			torrentFiles,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag: field.GormTag{
					"foreignKey": []string{"InfoHash"},
				},
			},
		),
		gen.FieldRelate(
			field.HasOne,
			"Pieces",
			torrentPieces,
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": []string{"InfoHash"},
				},
				JSONTag: "-",
			},
		),
		gen.FieldRelate(
			field.HasMany,
			"Tags",
			torrentTags,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag: field.GormTag{
					"foreignKey": []string{"InfoHash"},
				},
			},
		),
		infoHashType,
		infoHashReadOnly,
		gen.FieldType("files_status", "FilesStatus"),
		gen.FieldGORMTag("files_status", func(tag field.GormTag) field.GormTag {
			tag.Remove("default")
			return tag
		}),
		gen.FieldType("files_count", "NullUint"),
		gen.FieldGenType("extension", "String"),
		gen.FieldGORMTag("extension", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "false")
			return tag
		}),
		gen.FieldType("size", "uint"),
		gen.FieldIgnore("tsv"),
		createdAtReadOnly,
	)
	metadataSources := g.GenerateModel(
		"metadata_sources",
		readAndCreateField("key"),
		createdAtReadOnly,
	)
	contentCollections := g.GenerateModel(
		"content_collections",
		readAndCreateField("type"),
		gen.FieldRelate(
			field.BelongsTo,
			"MetadataSource",
			metadataSources,
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": {"Source"},
				},
			},
		),
		readAndCreateField("id"),
		createdAtReadOnly,
	)
	contentAttributes := g.GenerateModel(
		"content_attributes",
		readAndCreateField("content_type"),
		gen.FieldType("content_type", "ContentType"),
		readAndCreateField("content_source"),
		readAndCreateField("content_id"),
		gen.FieldRelate(
			field.BelongsTo,
			"MetadataSource",
			metadataSources,
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": {"Source"},
				},
			},
		),
		readAndCreateField("key"),
		createdAtReadOnly,
	)
	content := g.GenerateModel(
		"content",
		gen.FieldRelate(
			field.Many2Many,
			"Collections",
			contentCollections,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag: field.GormTag{
					"many2many": {"content_collections_content"},
				},
			},
		),
		gen.FieldRelate(
			field.HasMany,
			"Attributes",
			contentAttributes,
			&field.RelateConfig{
				RelateSlice: true,
			},
		),
		gen.FieldRelate(
			field.BelongsTo,
			"MetadataSource",
			metadataSources,
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": {"Source"},
				},
			},
		),
		readAndCreateField("type"),
		readAndCreateField("source"),
		readAndCreateField("id"),
		gen.FieldType("type", "ContentType"),
		gen.FieldGenType("type", "String"),
		gen.FieldType("release_date", "Date"),
		gen.FieldGenType("release_date", "Time"),
		gen.FieldType("release_year", "Year"),
		gen.FieldGenType("release_year", "Uint16"),
		gen.FieldType("original_language", "NullLanguage"),
		gen.FieldType("popularity", "NullFloat32"),
		gen.FieldType("vote_average", "NullFloat32"),
		gen.FieldType("vote_count", "NullUint"),
		gen.FieldType("runtime", "NullUint16"),
		gen.FieldType("adult", "NullBool"),
		gen.FieldType("tsv", "fts.Tsvector"),
		createdAtReadOnly,
	)
	contentCollectionContent := g.GenerateModelAs(
		"content_collections_content",
		"ContentCollectionContent",
		gen.FieldRelate(
			field.BelongsTo,
			"Content",
			content,
			&field.RelateConfig{},
		),
		gen.FieldRelate(
			field.BelongsTo,
			"Collection",
			contentCollections,
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": {"ContentCollectionType,ContentCollectionSource,ContentCollectionID"},
					"references": {"Type,Source,ID"},
				},
			},
		),
		gen.FieldType("content_type", "ContentType"),
	)
	torrentContentBaseOptions := []gen.ModelOpt{
		infoHashType,
		infoHashReadOnly,
		gen.FieldGenType("content_type", "String"),
		gen.FieldType("content_source", "NullString"),
		gen.FieldGenType("content_source", "String"),
		gen.FieldType("content_id", "NullString"),
		gen.FieldGenType("content_id", "String"),
		gen.FieldType("languages", "Languages"),
		gen.FieldGORMTag("languages", func(tag field.GormTag) field.GormTag {
			tag.Set("serializer", "json")
			return tag
		}),
		gen.FieldType("episodes", "Episodes"),
		gen.FieldGORMTag("episodes", func(tag field.GormTag) field.GormTag {
			tag.Set("serializer", "json")
			return tag
		}),
		gen.FieldType("video_resolution", "NullVideoResolution"),
		gen.FieldType("video_source", "NullVideoSource"),
		gen.FieldType("video_codec", "NullVideoCodec"),
		gen.FieldType("video_3d", "NullVideo3D"),
		gen.FieldType("video_modifier", "NullVideoModifier"),
		gen.FieldType("tsv", "fts.Tsvector"),
		createdAtReadOnly,
	}
	torrentContent := g.GenerateModel(
		"torrent_contents",
		append(
			[]gen.ModelOpt{
				gen.FieldGORMTag("id", func(tag field.GormTag) field.GormTag {
					tag.Set("<-", "false")
					return tag
				}),
				gen.FieldRelate(
					field.BelongsTo,
					"Torrent",
					torrents,
					&field.RelateConfig{
						GORMTag: field.GormTag{
							"foreignKey": {"InfoHash"},
							"references": {"InfoHash"},
						},
					},
				),
				gen.FieldRelate(
					field.BelongsTo,
					"Content",
					content,
					&field.RelateConfig{
						GORMTag: field.GormTag{
							"foreignKey": []string{"ContentType,ContentSource,ContentID"},
							"references": []string{"Type,Source,ID"},
						},
					},
				),
				gen.FieldType("content_type", "NullContentType"),
				gen.FieldType("seeders", "NullUint"),
				gen.FieldType("leechers", "NullUint"),
				gen.FieldType("size", "uint"),
				gen.FieldType("files_count", "NullUint"),
			},
			torrentContentBaseOptions...,
		)...,
	)
	torrentHints := g.GenerateModel(
		"torrent_hints",
		append(
			[]gen.ModelOpt{
				gen.FieldType("content_type", "ContentType"),
				gen.FieldType("release_year", "Year"),
			},
			torrentContentBaseOptions...,
		)...,
	)
	bloomFilters := g.GenerateModel(
		"bloom_filters",
		gen.FieldRename("bytes", "Filter"),
		gen.FieldType("bytes", "bloom.StableBloomFilter"),
		createdAtReadOnly,
	)
	keyValues := g.GenerateModel(
		"key_values",
		createdAtReadOnly,
	)
	queueJobs := g.GenerateModel(
		"queue_jobs",
		gen.FieldGORMTag("id", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "false")
			return tag
		}),
		gen.FieldGORMTag("queue", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		gen.FieldGORMTag("fingerprint", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		gen.FieldGORMTag("payload", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		gen.FieldType("status", "QueueJobStatus"),
		gen.FieldGenType("status", "String"),
		gen.FieldType("priority", "int"),
		gen.FieldType("retries", "uint"),
		gen.FieldType("max_retries", "uint"),
		gen.FieldGORMTag("max_retries", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		gen.FieldType("ran_at", "sql.NullTime"),
		gen.FieldType("deadline", "sql.NullTime"),
		gen.FieldGORMTag("deadline", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			return tag
		}),
		gen.FieldType("archival_duration", "Duration"),
		gen.FieldGORMTag("archival_duration", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "create")
			tag.Remove("default")
			return tag
		}),
		createdAtReadOnly,
	)

	g.ApplyBasic(
		torrentSources,
		torrentFiles,
		torrentsTorrentSources,
		torrentPieces,
		torrentTags,
		torrents,
		metadataSources,
		torrentContent,
		torrentHints,
		contentCollections,
		content,
		contentCollectionContent,
		contentAttributes,
		bloomFilters,
		keyValues,
		queueJobs,
	)

	return g
}

func readAndCreateField(columnName string) gen.ModelOpt {
	return gen.FieldGORMTag(columnName, func(tag field.GormTag) field.GormTag {
		tag.Set("<-", "create")
		return tag
	})
}
