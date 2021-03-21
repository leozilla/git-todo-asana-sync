package smarter.ecommerce.retailhub.commons.lang.util

import scala.util.chaining.scalaUtilChainingOps
import scala.util.{Failure, Try}

class RichTry[+T](private val original: Try[T]) extends AnyVal {

  def mapFailure(f: Throwable => Throwable): Try[T] = original.recoverWith {
    case err => Failure(f(err))
  }

  /** Use for performing side effects (like logging the exception) when this [[Try]] is a [[Failure]] */
  def tapFailure(effectF: Throwable => Unit): Try[T] = original.tap { t =>
    t.recover {
      case failure => effectF(failure)
    }
  }
}

object RichTry {
  implicit def tryToRichTry[T](original: Try[T]): RichTry[T] = new RichTry(original)
}