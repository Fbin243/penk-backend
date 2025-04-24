package rpc

import "github.com/jinzhu/copier"

func Map[From, To any](from *From, converters []copier.TypeConverter) (*To, error) {
	to := new(To)
	err := copier.CopyWithOption(to, from, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters:  converters,
	})

	return to, err
}

func MapSlice[From, To any](froms []From, converters []copier.TypeConverter) ([]To, error) {
	tos := make([]To, len(froms))
	err := copier.CopyWithOption(&tos, &froms, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters:  converters,
	})

	return tos, err
}
