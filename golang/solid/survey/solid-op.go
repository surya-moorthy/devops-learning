package survey

// import (
// 	"fmt"
// )

func ExportSurvey(s *Survey,exporter Exporter) error {

	// here there will be a modification if we try to add a new destination ,
	//  so we create a export interface and extended it based on our need.

	// switch dst {
	// case "s3":
	// 	// export to aws s3
	// 	return nil

	// case "gcs":
	// 	// export to gcs
	// 	return nil

	// default:
    //    return  fmt.Errorf(" unsupported destination: %s",dst)
	// }

	return exporter.Export(s)
}

type Exporter interface {
	Export(s *Survey) error
}

type S3Exporter struct {}

func(s3 *S3Exporter) Export(s *Survey) error {
	return nil
}

type GCSExporter struct {}

func(gcs *GCSExporter) Export(s *Survey) error {
	return nil
}

