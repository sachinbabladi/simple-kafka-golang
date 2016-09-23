package sarama

type ApiVersionsResponseBlock struct {
	ApiKey     int16
	MinVersion int16
	MaxVersion int16
}

func (r *ApiVersionsResponseBlock) encode(pe packetEncoder) error {
	pe.putInt16(r.ApiKey)
	pe.putInt16(r.MinVersion)
	pe.putInt16(r.MaxVersion)
	return nil
}

func (r *ApiVersionsResponseBlock) decode(pd packetDecoder) error {
	var err error

	if r.ApiKey, err = pd.getInt16(); err != nil {
		return err
	}

	if r.MinVersion, err = pd.getInt16(); err != nil {
		return err
	}

	if r.MaxVersion, err = pd.getInt16(); err != nil {
		return err
	}

	return nil
}

type ApiVersionsResponse struct {
	Err         KError
	ApiVersions []*ApiVersionsResponseBlock
}

func (r *ApiVersionsResponse) encode(pe packetEncoder) error {
	pe.putInt16(int16(r.Err))
	if err := pe.putArrayLength(len(r.ApiVersions)); err != nil {
		return err
	}
	for _, apiVersion := range r.ApiVersions {
		if err := apiVersion.encode(pe); err != nil {
			return err
		}
	}
	return nil
}

func (r *ApiVersionsResponse) decode(pd packetDecoder, version int16) error {
	if kerr, err := pd.getInt16(); err != nil {
		return err
	} else {
		r.Err = KError(kerr)
	}

	numBlocks, err := pd.getArrayLength()
	if err != nil {
		return err
	}

	r.ApiVersions = make([]*ApiVersionsResponseBlock, numBlocks)
	for i := 0; i < numBlocks; i++ {
		block := new(ApiVersionsResponseBlock)
		if err := block.decode(pd); err != nil {
			return err
		}
		r.ApiVersions[i] = block
	}

	return nil
}

func (r *ApiVersionsResponse) key() int16 {
	return 18
}

func (r *ApiVersionsResponse) version() int16 {
	return 0
}
