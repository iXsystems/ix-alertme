# $FreeBSD$

PORTNAME=	ix-alertme
PORTVERSION=	%%CHGVERSION%%
CATEGORIES=	misc

MAINTAINER=	ken@ixsystems.com
COMMENT=	Plugin-based alert notification system

LICENSE=	BSD2CLAUSE

USES=	go

USE_GITHUB=	yes
GH_ACCOUNT=	ixsystems
GH_PROJECT=	ix-alertme
GH_TAGNAME=	%%GHTAG%%

GH_TUPLE=		mholt:archiver:v3.2.0:archiver/src/github.com/mholt/archiver \
				andybalholm:brotli:71eb68c:brotli/src/github.com/andybalholm/brotli \
				golang:gddo:af0f2af:gddo/src/github.com/golang/gddo \
				klauspost:cpuid:v1.2.1:cpuid/src/github.com/klauspost/cpuid \
				dsnet:compress:v0.0.1:compress/src/github.com/dsnet/compress \
				golang:snappy:v0.0.1:snappy/src/github.com/golang/snappy \
				klauspost:compress:v1.7.3:klscompress/src/github.com/klauspost/compress \
				klauspost:pgzip:v1.2.1:klspost/src/github.com/klauspost/pgzip \
				nwaples:rardecode:v1.0.0:rardecode/src/github.com/nwaples/rardecode \
				pierrec:lz4:v2.2.4:lz4/src/github.com/pierrec/lz4 \
				ulikunitz:xz:v0.5.6:xz/src/github.com/ulikunitz/xz \
				xi2:xz:48954b6:xxz/src/github.com/xi2/xz \
				alecthomas:kingpin:v2.2.6:kingpin/src/gopkg.in/alecthomas/kingpin.v2 \
				alecthomas:template:a0175ee:template/src/github.com/alecthomas/template \
				alecthomas:units:c3de453:units/src/github.com/alecthomas/units
				

OPTIONS_DEFINE=	SMTPMAIL
OPTIONS_SUB=		yes
SMTPMAIL_DESC=	Bundle the smtp-email plugin into the package

.include <bsd.port.options.mk>

do-build:
	(cd ${WRKSRC}/src-go/ix-alertme ; \
		${GO_ENV} GOPATH=${WRKSRC} ${GO_CMD} build -o ix-alertme)

do-install:
	${INSTALL_PROGRAM} ${WRKSRC}/src-go/ix-alertme/ix-alertme ${STAGEDIR}${PREFIX}/bin/

post-install:
.if ${PORT_OPTIONS:MSMTPMAIL}
	${MKDIR} ${STAGEDIR}/${PREFIX}/${portname}/plugins/smtp-mail
	(cd ${WRKSRC}/provider-plugins/smtp-email &&
	${COPYTREE_BIN} smtp-email ${STAGEDIR}/${PREFIX}/${portname}/plugins/smtp-mail/smtp-mail)
	(cd ${WRKSRC}/provider-plugins/smtp-email &&
	${COPYTREE_SHARE} manifest.json ${STAGEDIR}/${PREFIX}/${portname}/plugins/smtp-mail/manifest.json)
.endif

.include <bsd.port.mk>
