package smarter.ecommerce.retailhub.commons.io.resources

import smarter.ecommerce.retailhub.common.io.GcsResource
import org.apache.commons.lang3.StringUtils

object GcsResourceConversion {

  // TODO: Use Java.URI for parsing?
  def urlToGcsResource(resourceUrl: String): Option[GcsResource] = {
    val bucket = StringUtils.substringBetween(resourceUrl, "gs://", "/")
    val path = StringUtils.substringAfter(resourceUrl, s"$bucket/")

    (bucket, path) match {
      case (bucket: String, path: String) => Some(GcsResource(bucket, path))
      case _ => None
    }
  }
}
