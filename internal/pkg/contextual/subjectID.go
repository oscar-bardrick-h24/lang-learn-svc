package contextual

import "context"

type subjectID struct{}

func SetSubjectID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, subjectID{}, id)
}

func GetSubjectID(ctx context.Context) string {
	return ctx.Value(subjectID{}).(string)
}
