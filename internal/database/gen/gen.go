package gen

import (
	"github.com/iancoleman/strcase"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"path"
	"runtime"
)

func readAndCreateField(columnName string) gen.ModelOpt {
	return gen.FieldGORMTag(columnName, func(tag field.GormTag) field.GormTag {
		tag.Set("<-", "create")
		return tag
	})
}

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
		gen.FieldType("size", "uint64"),
		gen.FieldType("index", "uint32"),
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
	torrentsTorrentSources := g.GenerateModel(
		"torrents_torrent_sources",
		readAndCreateField("source"),
		infoHashType,
		infoHashReadOnly,
		gen.FieldType("seeders", "NullUint"),
		gen.FieldType("leechers", "NullUint"),
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
		gen.FieldGenType("extension", "String"),
		gen.FieldGORMTag("extension", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "false")
			return tag
		}),
		gen.FieldType("size", "uint64"),
		gen.FieldType("piece_length", "NullUint64"),
		gen.FieldJSONTag("pieces", "-"),
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
		gen.FieldType("release_date", "Date"),
		gen.FieldGenType("release_date", "Time"),
		gen.FieldType("release_year", "Year"),
		gen.FieldType("original_language", "NullLanguage"),
		gen.FieldType("popularity", "NullFloat32"),
		gen.FieldType("vote_average", "NullFloat32"),
		gen.FieldType("vote_count", "NullUint"),
		gen.FieldType("runtime", "NullUint16"),
		gen.FieldType("adult", "NullBool"),
		gen.FieldIgnore("tsv"),
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
	torrentContent := g.GenerateModel(
		"torrent_contents",
		gen.FieldRelate(
			field.BelongsTo,
			"Torrent",
			torrents,
			&field.RelateConfig{
				GORMTag: field.GormTag{
					"foreignKey": {"InfoHash"},
				},
			},
		),
		gen.FieldRelate(
			field.BelongsTo,
			"Content",
			content,
			&field.RelateConfig{},
		),
		gen.FieldGORMTag("id", func(tag field.GormTag) field.GormTag {
			tag.Set("<-", "false")
			return tag
		}),
		infoHashType,
		infoHashReadOnly,
		gen.FieldType("content_type", "NullContentType"),
		gen.FieldType("content_source", "NullString"),
		gen.FieldType("content_id", "NullString"),
		gen.FieldType("release_date", "Date"),
		gen.FieldGenType("release_date", "Time"),
		gen.FieldType("release_year", "Year"),
		gen.FieldType("external_ids", "maps.StringMap[string]"),
		gen.FieldGORMTag("external_ids", func(tag field.GormTag) field.GormTag {
			tag.Set("serializer", "json")
			return tag
		}),
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
		gen.FieldType("video_3d", "NullVideo3d"),
		gen.FieldType("video_modifier", "NullVideoModifier"),
		gen.FieldIgnore("tsv"),
		createdAtReadOnly,
	)
	bloomFilters := g.GenerateModel(
		"bloom_filters",
		gen.FieldRename("bytes", "Filter"),
		gen.FieldType("bytes", "bloom.StableBloomFilter"),
		createdAtReadOnly,
	)

	g.ApplyBasic(
		torrentSources,
		torrentFiles,
		torrentsTorrentSources,
		torrentTags,
		torrents,
		metadataSources,
		torrentContent,
		contentCollections,
		content,
		contentCollectionContent,
		contentAttributes,
		bloomFilters,
	)

	return g
}
