package mysql

import "apu/pkg/store/mysql/query"

type PipelineOptions struct {
	IsCounted    bool
	IsDetailed   bool
	IsStructured bool
	IsEmbedded   bool
	IsIndexed    bool
}

// UpdatePipeline 将笔记加入任务管道。
func UpdatePipeline(noteID uint64, opts ...PipelineOptions) error {
	opt := PipelineOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	p, err := query.Pipeline.Where(query.Pipeline.NoteID.Eq(noteID)).FirstOrCreate()
	if err != nil {
		return err
	}

	needSave := false

	if opt.IsCounted {
		p.IsCounted = true
		needSave = true
	}

	if opt.IsDetailed {
		p.IsDetailed = true
		needSave = true
	}
	if opt.IsStructured {
		p.IsStructured = true
		needSave = true
	}

	if opt.IsEmbedded {
		p.IsEmbedded = true
		needSave = true
	}
	if opt.IsIndexed {
		p.IsIndexed = true
		needSave = true
	}

	if needSave {
		err = query.Pipeline.Save(p)
		if err != nil {
			return err
		}
	}

	return nil
}
