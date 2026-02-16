package errors

func HandleError(err error) (int, string, string) {
    
    // STEP 1: Is it a Local Error from my 'pkg/apperrors'?
    // We check our "Dictionary" (the map we made).
    if resp, ok := errorMap[err]; ok {
        return resp.HTTPCode, resp.Slug, err.Error()
    }

    // STEP 2: If it's NOT in our dictionary, is it a Remote gRPC error?
    // We check if it can be converted to a gRPC Status.
    if st, ok := status.FromError(err); ok {
        // We use a helper to turn gRPC codes into HTTP
        return MapGrpcCodeToHttp(st.Code(), st.Message())
    }

    // STEP 3: We have no idea what this is.
    // Safety net to prevent the app from leaking sensitive details.
    return 500, "INTERNAL_SERVER_ERROR", "Something went wrong"
}