package smarter.ecommerce.retailhub.commons.io.resources

import akka.http.scaladsl.model.ContentTypes
import akka.stream.alpakka.googlecloud.storage.scaladsl.GCStorage
import akka.stream.scaladsl.{Keep, Sink}
import akka.util.ByteString
import com.google.common.io.ByteStreams
import org.scalatest.concurrent.IntegrationPatience
import org.scalatest.matchers.must.Matchers
import org.scalatest.wordspec.AnyWordSpec
import org.scalatest.{BeforeAndAfterAll, Suite}
import smarter.ecommerce.cloudflow.testutils.flows.FlowTestHelper
import smarter.ecommerce.retailhub.testutil.akka.ActorSystemSupport
import smarter.ecommerce.retailhub.testutil.gcs.GcsTestSuiteMixin
import smarter.ecommerce.retailhub.testutil.scalatest.concurrent.ScalaFutures
import smarter.ecommerce.retailhub.testutil.scalatest.concurrent.Timeouts.defaultTimeout

import java.time.{Duration, Instant}
import scala.collection.mutable.ListBuffer

// FIXME: centralize this config, already repeated in three places
trait TranscodingSupportExtSpecSupport extends GcsTestSuiteMixin {
  self: Suite =>

  private val stage = "dev" // TODO inject stage param from build, once we move to multi-stage environments
  override val gcloudProject = s"smec-retailhub-platform-$stage"
  override val bucket = s"$gcloudProject-xspec-ingestion"
  override val customCredentialsJson = None
}

class TranscodingSupportExtSpec extends AnyWordSpec
  with Matchers
  with BeforeAndAfterAll
  with ActorSystemSupport
  with TranscodingSupportExtSpecSupport
  with ScalaFutures
  with IntegrationPatience {

  private def readFile(fileName: String) = getClass.getClassLoader.getResourceAsStream(fileName)

  val gcsResourceAccessor = new GcsResourceAccessor()

  val gzippedFile = readFile("test.txt.gz")
  val gzippedFileContent = ByteString(ByteStreams.toByteArray(gzippedFile))
  val gunzippedFile = readFile("test.txt")
  val gunzippedFileContent = ByteString(ByteStreams.toByteArray(gunzippedFile))

  val storedObject = new BucketResource {
    override val bucket: String = gcsResource.bucket

    override val path: String = gcsResource.fullPath("test")
  }

  override def beforeAll(): Unit = {
    super.beforeAll()

    implicit val ec = actorSystem.dispatcher

    actorSystem.log.info(s"Using base path ${gcsResource.basePath}")

    gcs.uploadObject(storedObject.path, gzippedFileContent, ContentTypes.`text/plain(UTF-8)`, Some("gzip")).futureValue
  }

  override def afterAll(): Unit = {
    gcs.deleteAllObjectsInFolder(gcsResource.basePath)
  }

  "successfully download compressed file from GCS" in {
    val resultBytes = FlowTestHelper(gcsResourceAccessor.readByteString)
      .fromSingle(storedObject)
      .runWithSeq(defaultTimeout).reduce(_ ++ _)

    resultBytes mustEqual gunzippedFileContent
  }

  "download speed benchmark" ignore {
    val FilesNamesPattern = (1 to 5).map(n => s"$n")

    val testRepetition = 20
    val baseDir = "zip-test/performance-test"

    def downloadAndRecordTime(fileName: String): Duration = {
      actorSystem.log.info(s"downloading file $fileName")

      val start = Instant.now()
      GCStorage
        .download(bucket, s"$baseDir/$fileName")
        .collect { case Some(src) => src }
        .toMat(Sink.ignore)(Keep.right)
        .run()
        .futureValue

      val end = Instant.now()

      Duration.between(start, end)
    }

    def getObjectSize(fileName: String): Double = {
      val sObject = GCStorage
        .getObject(bucket, s"$baseDir/$fileName")
        .collect {
          case Some(src) =>
            src
        }.toMat(Sink.head)(Keep.right)
        .run()
        .futureValue

      sObject.size / 1000000.00 // size is in bytes and then converted to megabytes; 1 MB = 1000000 bytes
    }

    val downloadSpeeds = ListBuffer.empty[(List[Long], List[Long])]

    "be calculated" in {

      // warmup, file 6 is out of the benchmark
      downloadAndRecordTime("6")
      downloadAndRecordTime("6.gz")

      FilesNamesPattern.foreach(file => {
        val gunzippedFileSpeeds = ListBuffer.empty[Long]
        val gzippedFileSpeeds = ListBuffer.empty[Long]

        (1 to testRepetition).foreach(rep => {
          actorSystem.log.info(s"file $file, repetition $rep")

          gunzippedFileSpeeds += downloadAndRecordTime(file).toMillis
          gzippedFileSpeeds += downloadAndRecordTime(s"$file.gz").toMillis
        })
        val resultTuple = (gunzippedFileSpeeds.toList, gzippedFileSpeeds.toList)
        actorSystem.log.info(s"file $file, results: $resultTuple")

        downloadSpeeds += resultTuple
      })

      downloadSpeeds.toList.zipWithIndex.foreach {
        case ((gunzippedSpeeds, gzippedSpeeds), index) =>
          val file = index + 1

          val (min, max, mean) = (gunzippedSpeeds.min, gunzippedSpeeds.max, gunzippedSpeeds.sum / gunzippedSpeeds.size)
          val (minG, maxG, meanG) = (gzippedSpeeds.min, gunzippedSpeeds.max, gzippedSpeeds.sum / gzippedSpeeds.size)

          val gzippedSize = getObjectSize(s"$file.gz")
          val unzippedSize = getObjectSize(s"$file")

          actorSystem.log.info(s"result for file $index, gzipped size: ${gzippedSize}MB, unzipped size:${unzippedSize}MB")
          actorSystem.log.info(s"result for file $file, mean speed gzipped: ${meanG}ms, mean speed unzipped: ${mean}ms")
          actorSystem.log.info(s"result for file $file, min speed gzipped: ${minG}ms, min speed unzipped: ${min}ms")
          actorSystem.log.info(s"result for file $file, max speed gzipped: ${maxG}ms, max speed unzipped: ${max}ms")
      }
    }
  }

}
